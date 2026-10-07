package main

import (
	"bytes"
	"go/format"
	"io/fs"
	"strings"
	"testing"
)

func TestTemplatesRender(t *testing.T) {
	p := newProject("myservice", "", "1.23")

	paths, err := fs.Glob(templatesFS, templateRoot)
	if err != nil || len(paths) == 0 {
		t.Fatalf("templates not embedded: %v", err)
	}

	var rendered int

	err = fs.WalkDir(templatesFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if !strings.HasSuffix(path, templateExt) {
			t.Errorf("%s: template must end in %s or the Go tool will compile it", path, templateExt)
			return nil
		}

		content, err := render(path, p)
		if err != nil {
			t.Errorf("%s: %v", path, err)
			return nil
		}

		rendered++

		if !strings.HasSuffix(path, ".go"+templateExt) {
			return nil
		}

		src, err := format.Source(content)
		if err != nil {
			t.Errorf("%s: rendered output is not valid Go: %v", path, err)
			return nil
		}

		if !bytes.Equal(src, content) {
			t.Errorf("%s: rendered output is not gofmt-clean", path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if rendered == 0 {
		t.Fatal("no templates rendered")
	}
}

func TestNewProjectModule(t *testing.T) {
	if got := newProject("svc", "", "1.23").Module; got != modulePrefix+"svc" {
		t.Errorf("empty module = %q, want %q", got, modulePrefix+"svc")
	}

	if got := newProject("svc", "example.com/foo/bar", "1.23").Module; got != "example.com/foo/bar" {
		t.Errorf("explicit module = %q, want example.com/foo/bar", got)
	}
}

func TestOutputPath(t *testing.T) {
	tests := map[string]string{
		"templates/app/app.go.tmpl":                        "app/app.go",
		"templates/.gitignore.tmpl":                        ".gitignore",
		"templates/app/model/.gitkeep.tmpl":                "app/model/.gitkeep",
		"templates/go.mod.tmpl":                            "go.mod",
		"templates/specs/server/params/model/echo.go.tmpl": "specs/server/params/model/echo.go",
	}

	for in, want := range tests {
		if got := outputPath(in); got != want {
			t.Errorf("outputPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDotfileTemplatesAreEmbedded(t *testing.T) {
	want := []string{
		templateRoot + "/.gitignore.tmpl",
		templateRoot + "/app/helpers/.gitkeep.tmpl",
		templateRoot + "/app/model/.gitkeep.tmpl",
		templateRoot + "/app/repository/.gitkeep.tmpl",
		templateRoot + "/specs/client/.gitkeep.tmpl",
	}

	for _, path := range want {
		if _, err := templatesFS.ReadFile(path); err != nil {
			t.Errorf("%s not embedded (does the directive still say all:?): %v", path, err)
		}
	}
}
