package templates

import (
	"html/template"
	"path/filepath"
)

var fm = template.FuncMap{}

type TemplateCache struct {
	cache *template.Template
}

func NewTemplateCache(path string) *TemplateCache {
	//example path: "web/templates/"
	t := template.Must(template.New("*").Funcs(fm).ParseGlob(path + "*.html"))
	files, err := filepath.Glob(path + "*/*.html")
	if len(files) > 0 && err == nil {
		t = template.Must(t.ParseGlob(path + "*/*.html"))
	}
	return &TemplateCache{
		cache: template.Must(template.New("*").Funcs(fm).ParseGlob(path)),
	}
}
