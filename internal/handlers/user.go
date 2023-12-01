package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/cookie"
	pwd "github.com/zhang2092/mediahls/internal/pkg/password"
)

func (server *Server) registerView(w http.ResponseWriter, r *http.Request) {
	// 是否已经登录
	server.isRedirect(w, r)
	renderRegister(w, nil)
}

func renderRegister(w http.ResponseWriter, data any) {
	render(w, data, "web/templates/user/register.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
}

// type userResponse struct {
// 	Username          string    `json:"username"`
// 	FullName          string    `json:"full_name"`
// 	Email             string    `json:"email"`
// 	PasswordChangedAt time.Time `json:"password_changed_at"`
// 	CreatedAt         time.Time `json:"created_at"`
// }

// func newUserResponse(user db.User) userResponse {
// 	return userResponse{
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		CreatedAt: user.CreatedAt,
// 	}
// }

func viladatorRegister(email, username, password string) (*respErrs, bool) {
	ok := true
	errs := &respErrs{
		Email:    email,
		Username: username,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		errs.EmailErr = "请填写正确的邮箱地址"
		ok = false
	}
	if !ValidateRxUsername(username) {
		errs.UsernameErr = "名称(6-20,字母,数字)"
		ok = false
	}
	if !ValidatePassword(password) {
		errs.PasswordErr = "密码(8-20位)"
		ok = false
	}

	return errs, ok
}

type respErrs struct {
	Authorize
	Summary     string
	Email       string
	Username    string
	Password    string
	EmailErr    string
	UsernameErr string
	PasswordErr string
}

func (server *Server) register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	errs, ok := viladatorRegister(email, username, password)
	if !ok {
		renderRegister(w, errs)
		return
	}

	hashedPassword, err := pwd.BcryptHashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arg := db.CreateUserParams{
		ID:             genId(),
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
	}

	_, err = server.store.CreateUser(r.Context(), arg)
	if err != nil {
		if server.store.IsUniqueViolation(err) {
			errs.Summary = "邮箱或名称已经存在"
			renderRegister(w, errs)
			return
		}

		errs.Summary = "请求网络错误,请刷新重试"
		renderRegister(w, errs)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)

	// rsp := newUserResponse(user)
	// Respond(w, "ok", rsp, http.StatusOK)
}

func (server *Server) loginView(w http.ResponseWriter, r *http.Request) {
	// 是否已经登录
	server.isRedirect(w, r)
	renderLogin(w, nil)
}

func renderLogin(w http.ResponseWriter, data any) {
	render(w, data, "web/templates/user/login.html.tmpl", "web/templates/base/header.html.tmpl", "web/templates/base/footer.html.tmpl")
}

// type loginUserResponse struct {
// 	AccessToken          string       `json:"access_token"`
// 	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
// 	User                 userResponse `json:"user"`
// }

func viladatorLogin(email, password string) (*respErrs, bool) {
	ok := true
	errs := &respErrs{
		Email:    email,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		errs.EmailErr = "请填写正确的邮箱地址"
		ok = false
	}
	if len(password) == 0 {
		errs.PasswordErr = "请填写正确的密码"
		ok = false
	}

	return errs, ok
}

func (server *Server) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		renderLogin(w, respErrs{Summary: "请求网络错误,请刷新重试"})
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	errs, ok := viladatorLogin(email, password)
	if !ok {
		renderLogin(w, errs)
		return
	}

	ctx := r.Context()
	user, err := server.store.GetUserByEmail(ctx, email)
	if err != nil {
		if server.store.IsNoRows(sql.ErrNoRows) {
			errs.Summary = "邮箱或密码错误"
			renderLogin(w, errs)
			return
		}

		errs.Summary = "请求网络错误,请刷新重试"
		renderLogin(w, errs)
		return
	}

	err = pwd.BcryptComparePassword(user.HashedPassword, password)
	if err != nil {
		errs.Summary = "邮箱或密码错误"
		renderLogin(w, errs)
		return
	}

	encoded, err := server.secureCookie.Encode(AuthorizeCookie, &Authorize{ID: user.ID, Name: user.Username})
	if err != nil {
		errs.Summary = "请求网络错误,请刷新重试(cookie)"
		renderLogin(w, errs)
		return
	}

	c := cookie.NewCookie(cookie.AuthorizeName, encoded, time.Now().Add(time.Duration(7200)*time.Second))
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie.DeleteCookie(w, cookie.AuthorizeName)
	http.Redirect(w, r, "/login", http.StatusFound)
}
