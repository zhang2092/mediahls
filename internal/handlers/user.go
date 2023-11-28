package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/zhang2092/mediahls/internal/db"
	pwd "github.com/zhang2092/mediahls/internal/pkg/password"
)

func (server *Server) registerView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/user/register.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	hashedPassword, err := pwd.BcryptHashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arg := db.CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
	}

	user, err := server.store.CreateUser(r.Context(), arg)
	if err != nil {
		if server.store.IsUniqueViolation(err) {
			http.Error(w, "数据已经存在", http.StatusInternalServerError)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rsp := newUserResponse(user)
	Respond(w, "ok", rsp, http.StatusOK)
}

func (server *Server) loginView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/user/login.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type loginUserResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 userResponse `json:"user"`
}

func (server *Server) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	ctx := r.Context()

	user, err := server.store.GetUserByName(ctx, username)
	if err != nil {
		if server.store.IsNoRows(sql.ErrNoRows) {
			http.Error(w, "用户不存在", http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pwd.BcryptComparePassword(password, user.HashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		server.conf.AccessTokenDuration,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rsp := loginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
		User:                 newUserResponse(user),
	}
	Respond(w, "ok", rsp, http.StatusOK)
}
