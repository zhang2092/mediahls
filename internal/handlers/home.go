package handlers

import (
	"html/template"
	"net/http"
)

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/home.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
