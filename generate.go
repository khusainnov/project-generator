package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

var emptyDirs = []string{
	"app/helpers",
	"app/model",
	"app/repository",
	"specs/client",
}

type Project struct {
	Name      string
	Module    string
	EnvPrefix string
	GoVersion string
}

func newProject(name, goVer string) Project {
	return Project{
		Name:      name,
		Module:    modulePrefix + name,
		EnvPrefix: strings.ToUpper(name),
		GoVersion: goVer,
	}
}

func generateProjectStructure(p Project) error {
	for _, dir := range emptyDirs {
		err := os.MkdirAll(filepath.Join(p.Name, dir), 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}

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

	return buf.Bytes(), nil
}

func outputPath(path string) string {
	rel := strings.TrimSuffix(strings.TrimPrefix(path, templateRoot+"/"), templateExt)

	if base := filepath.Base(rel); base == "gitignore" {
		rel = filepath.Join(filepath.Dir(rel), "."+base)
	}

	return rel
}
