package handler

import (
	"html/template"
	"net/http"
	"regexp"

	"github.com/KushamiNeko/go_fun/utils/web"
)

var templates *template.Template

type TemplateData struct {
	ID    string
	Class string
	Data  interface{}
}

func NewTemplateData(id, cls string, data interface{}) *TemplateData {
	return &TemplateData{
		ID:    id,
		Class: cls,
		Data:  data,
	}
}

func init() {
	var err error

	templates, err = template.New("").Funcs(
		template.FuncMap{
			"TemplateData": NewTemplateData,
		},
	).ParseGlob("market_wizards/templates/**/**/*.html")
	if err != nil {
		panic(err)
	}
}

type ViewHandler struct{}

func (v *ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		const pattern = `/view/(practice)/*.*`

		regex := regexp.MustCompile(pattern)
		match := regex.FindAllStringSubmatch(r.RequestURI, -1)
		if match == nil {
			http.NotFound(w, r)
			return
		}

		switch match[0][1] {
		case "practice":
			v.practice(w, r)
		default:
			http.NotFound(w, r)
			return
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (v *ViewHandler) practice(w http.ResponseWriter, r *http.Request) {
	web.WriteTemplate(w, templates, "Practice", nil)
}
