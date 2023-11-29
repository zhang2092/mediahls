package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/zhang2092/mediahls/internal/pkg/convert"
)

func (server *Server) authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := server.withCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		b, err := json.Marshal(u)
		if err != nil {
			log.Printf("json marshal authorize user: %v", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUser, b)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) withCookie(r *http.Request) (*authorize, error) {
	cookie, err := r.Cookie(AuthorizeCookie)
	if err != nil {
		log.Printf("get cookie: %v", err)
		return nil, err
	}

	u := &authorize{}
	err = server.secureCookie.Decode(AuthorizeCookie, cookie.Value, u)
	if err != nil {
		log.Printf("secure decode cookie: %v", err)
		return nil, err
	}

	return u, nil
}

func withUser(ctx context.Context) *authorize {
	var result authorize
	ctxValue, err := convert.ToByteE(ctx.Value(ContextUser))
	log.Printf("ctx: %s", ctxValue)
	if err != nil {
		log.Printf("1: %v", err)
		return nil
	}
	err = json.Unmarshal(ctxValue, &result)
	if err != nil {
		log.Printf("2: %v", err)
		return nil
	}

	return &result
}
