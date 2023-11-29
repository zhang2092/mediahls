package handlers

import (
	"net/http"
)

type pageData struct {
	AuthID   string
	AuthName string
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	pd := pageData{}
	auth, err := server.withCookie(r)
	if err == nil {
		pd.AuthID = auth.AuthID
		pd.AuthName = auth.AuthName
	}
	renderHome(w, pd)
}

func renderHome(w http.ResponseWriter, data any) {
	render(w, data, "web/templates/home.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
}
