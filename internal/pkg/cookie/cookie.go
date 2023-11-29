package cookie

import (
	"net/http"
	"time"
)

const (
	AuthorizeName = "authorize"
)

func NewCookie(name, value string, expired time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   false, // true->只能https站点操作
		HttpOnly: true,  // true->js不能捕获
		Expires:  expired,
	}
}

func SetCookie(w http.ResponseWriter, name, value string, expired time.Time) {
	cookie := NewCookie(name, value, expired)
	http.SetCookie(w, cookie)
}

func ReadCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := NewCookie(name, "", time.Now().Add(time.Duration(-10)*time.Second))
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
