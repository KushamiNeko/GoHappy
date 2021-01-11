package handler

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/KushamiNeko/go_fun/utils/web"
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

	w.Write(web.CleanAll(buffer.Bytes()))
}
