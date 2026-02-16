package web

import (
	"bytes"
	"html/template"
	"net/http"
)

type Renderer struct {
	Templ *template.Template
}

func NewRenderer(globPattern string) (*Renderer, error) {
	t, err := template.ParseGlob(globPattern)
	if err != nil {
		return nil, err
	}
	return &Renderer{
		Templ: t,
	}, nil
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, status int, templateName string, data any) error {
	
	var contentBuffer bytes.Buffer
	err := r.Templ.ExecuteTemplate(&contentBuffer, templateName, data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(contentBuffer.Bytes())
	return nil
}
