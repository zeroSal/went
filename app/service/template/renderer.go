package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	tmpl "text/template"
)

type renderer struct {
	templates fs.FS
	results   map[string][]byte
}

var _ Renderer = (*renderer)(nil)

func NewRenderer(templates fs.FS) Renderer {
	return &renderer{
		templates: templates,
		results:   make(map[string][]byte),
	}
}

func (r *renderer) Render(data map[string]string) error {
	templatesDir := r.templates
	_, err := templatesDir.Open("templates")
	if err != nil {
		templatesDir, err = fs.Sub(r.templates, "embed")
		if err != nil {
			return fmt.Errorf("get embed subdirectory: %w", err)
		}
	}

	err = fs.WalkDir(templatesDir, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		t, err := tmpl.ParseFS(templatesDir, p)
		if err != nil {
			return fmt.Errorf("parse template %s: %w", p, err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute template %s: %w", p, err)
		}

		resultPath := renderPath(p, data)
		r.results[resultPath] = buf.Bytes()
		return nil
	})

	return err
}

func (r *renderer) GetFiles() []string {
	files := make([]string, 0, len(r.results))
	for f := range r.results {
		files = append(files, f)
	}
	return files
}

func (r *renderer) Read(path string) ([]byte, bool) {
	c, ok := r.results[path]
	return c, ok
}

func renderPath(p string, _ map[string]string) string {
	p = strings.TrimPrefix(p, "templates/")
	p = strings.TrimSuffix(p, ".template")
	return p
}
