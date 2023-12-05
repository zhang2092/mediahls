package handlers

import (
	"net/http"

	"github.com/rs/xid"
)

const (
	AuthorizeCookie        = "authorize"
	ContextUser     ctxKey = "context_user"
)

type ctxKey string

type Authorize struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func genId() string {
	id := xid.New()
	return id.String()
}

func (server *Server) isRedirect(w http.ResponseWriter, r *http.Request) {
	u := withUser(r.Context())
	if u != nil {
		// 已经登录, 直接到首页
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
