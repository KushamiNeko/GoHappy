package handler

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/KushamiNeko/go_utils/cleaner"
)

func renderView(w http.ResponseWriter, viewTemplates string, data interface{}) {
	buffer := bytes.Buffer{}

	index := template.Must(template.New("Index").ParseGlob("templates/views/index/*.html"))
	index = template.Must(index.ParseGlob("templates/_components/**/*.html"))

	index = template.Must(index.ParseGlob(viewTemplates))

	err := index.ExecuteTemplate(&buffer, "Index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(cleaner.Clean(buffer.Bytes()))
}

func writeTemplate(w http.ResponseWriter, temp *template.Template, name string, data interface{}, cb func()) {
	buffer := bytes.Buffer{}

	err := temp.ExecuteTemplate(
		&buffer,
		name,
		data,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cb != nil {
		cb()
	}

	w.Write(cleaner.Clean(buffer.Bytes()))
}
