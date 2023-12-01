package handlers

import (
	"html/template"
	"net/http"

	"github.com/zhang2092/mediahls/internal/pkg/logger"
)

// func render(w http.ResponseWriter, data any, tmpls ...string) {
// 	t, err := template.ParseFiles(tmpls...)
// 	if err != nil {
// 		log.Printf("template parse: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	err = t.Execute(w, data)
// 	if err != nil {
// 		log.Printf("template execute: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }

// renderLayout 渲染方法 带框架
func renderLayout(w http.ResponseWriter, data any, tmpls ...string) {
	tmpls = append(tmpls, "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
	t, err := template.ParseFiles(tmpls...)
	if err != nil {
		logger.Logger.Errorf("template parse: %v, %v", tmpls, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		logger.Logger.Errorf("template execute: %v, %v", tmpls, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
