package handler

import (
	"net/http"
	"regexp"
)

type ViewHandler struct{}

func (v *ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		const pattern = `/view/(study|practice)/*.*`

		regex := regexp.MustCompile(pattern)
		if !regex.MatchString(r.RequestURI) {
			http.NotFound(w, r)
			return
		}

		match := regex.FindAllStringSubmatch(r.RequestURI, -1)

		switch match[0][1] {
		case "study":
			v.getStudy(w, r)
		case "practice":
			v.getPractice(w, r)
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (v *ViewHandler) getStudy(w http.ResponseWriter, r *http.Request) {
	renderView(w, "templates/views/study/*.html", nil)
}

func (v *ViewHandler) getPractice(w http.ResponseWriter, r *http.Request) {
	renderView(w, "templates/views/practice/*.html", nil)
}
