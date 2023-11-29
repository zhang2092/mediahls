package handlers

import (
	"net/http"

	"github.com/rs/xid"
	"github.com/zhang2092/mediahls/internal/pkg/cookie"
)

const (
	AuthorizeCookie             = "authorize"
	ContextUser     CtxTypeUser = "context_user"
)

type CtxTypeUser string

type authorize struct {
	AuthID   string `json:"auth_id"`
	AuthName string `json:"auth_name"`
}

func genId() string {
	id := xid.New()
	return id.String()
}

func (server *Server) isRedirect(w http.ResponseWriter, r *http.Request) {
	_, err := server.withCookie(r)
	if err != nil {
		// 1. 删除cookie
		cookie.DeleteCookie(w, cookie.AuthorizeName)
		return
	}

	// cookie 校验成功
	http.Redirect(w, r, "/", http.StatusFound)
}
