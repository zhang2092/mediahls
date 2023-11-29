package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func render(w http.ResponseWriter, data any, tmpls ...string) {
	t, err := template.ParseFiles(tmpls...)
	if err != nil {
		log.Printf("template parse: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("template execute: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
