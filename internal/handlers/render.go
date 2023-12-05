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
		"currentUser": func() *Authorize {
			return withUser(r.Context())
		},
	})

	tpl := template.Must(t.Clone())

	// compress
	// "github.com/tdewolff/minify/v2"
	// "github.com/tdewolff/minify/v2/css"
	// "github.com/tdewolff/minify/v2/html"
	// "github.com/tdewolff/minify/v2/js"
	// m := minify.New()
	// m.AddFunc("text/css", css.Minify)
	// m.AddFunc("text/html", html.Minify)
	// m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	// pages := []string{
	// 	tmpl,
	// 	"base/header.html.tmpl",
	// 	"base/footer.html.tmpl",
	// }
	// for _, page := range pages {
	// 	b, err := fs.ReadFile(server.templateFS, page)
	// 	if err != nil {
	// 		logger.Logger.Errorf("fs read file: %s, %v", page, err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	mb, err := m.Bytes("text/html", b)
	// 	if err != nil {
	// 		logger.Logger.Errorf("minify bytes: %s, %v", page, err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	tpl, err = tpl.Parse(string(mb))
	// 	if err != nil {
	// 		logger.Logger.Errorf("template parse: %s, %v", page, err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	tpl, err := tpl.ParseFS(server.templateFS, tmpl, "base/header.html.tmpl", "base/footer.html.tmpl")
	if err != nil {
		logger.Logger.Errorf("template parse: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, data); err != nil {
		logger.Logger.Errorf("template execute: %s, %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
