package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
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
func (server *Server) renderLayout(w http.ResponseWriter, r *http.Request, data any, tmpl string) {
	t := template.New(filepath.Base(tmpl))
	t = t.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})

	tpl := template.Must(t.Clone())
	tpl, err := tpl.ParseFS(server.templateFS, tmpl, "base/header.html.tmpl", "base/footer.html.tmpl")
	if err != nil {
		logger.Logger.Errorf("template parse: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		logger.Logger.Errorf("template execute: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
