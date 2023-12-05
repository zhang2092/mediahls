package handlers

import (
	"context"
	"net/http"
)

func (server *Server) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := withUser(r.Context())
		if u == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (server *Server) setUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthorizeCookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		u := Authorize{}
		err = server.secureCookie.Decode(AuthorizeCookie, cookie.Value, &u)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withUser(ctx context.Context) *Authorize {
	val := ctx.Value(ContextUser)
	if u, ok := val.(Authorize); ok {
		return &u
	}

	return nil
}
