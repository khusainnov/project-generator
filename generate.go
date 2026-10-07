package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed all:templates
var templatesFS embed.FS

const (
	templateRoot = "templates"
	templateExt  = ".tmpl"
	modulePrefix = "github.com/khusainnov/"
)

type Project struct {
	Name      string
	Module    string
	EnvPrefix string
	GoVersion string
}

func newProject(name, module, goVer string) Project {
	if module == "" {
		module = modulePrefix + name
	}

	return Project{
		Name:      name,
		Module:    module,
		EnvPrefix: envPrefix(name),
		GoVersion: goVer,
	}
}

var notEnvNameChar = regexp.MustCompile(`[^A-Z0-9_]`)

func envPrefix(name string) string {
	s := notEnvNameChar.ReplaceAllString(strings.ToUpper(name), "_")

	if s == "" || (s[0] >= '0' && s[0] <= '9') {
		return "_" + s
	}

	return s
}

func generateProjectStructure(p Project) error {
	return fs.WalkDir(templatesFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		rel := outputPath(path)

		content, err := render(path, p)
		if err != nil {
			return err
		}

		out := filepath.Join(p.Name, rel)
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return fmt.Errorf("failed to create dir for %s: %w", rel, err)
		}

		if err := os.WriteFile(out, content, 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", rel, err)
		}

		return nil
	})
}

func render(path string, p Project) ([]byte, error) {
	tmpl, err := template.New(filepath.Base(path)).ParseFS(templatesFS, path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return nil, fmt.Errorf("failed to render template %s: %w", path, err)
	}

	if !strings.HasSuffix(path, ".go"+templateExt) {
		return buf.Bytes(), nil
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("template %s produced invalid Go: %w", path, err)
	}

	return src, nil
}

func outputPath(path string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, templateRoot+"/"), templateExt)
}
