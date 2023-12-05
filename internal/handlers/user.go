package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/zhang2092/mediahls/internal/db"
	"github.com/zhang2092/mediahls/internal/pkg/cookie"
	pwd "github.com/zhang2092/mediahls/internal/pkg/password"
)

// obj

// registerPageData 注册页面数据
type registerPageData struct {
	Summary     string
	Email       string
	EmailMsg    string
	Username    string
	UsernameMsg string
	Password    string
	PasswordMsg string
}

// loginPageData 登录页面数据
type loginPageData struct {
	Summary     string
	Email       string
	EmailMsg    string
	Password    string
	PasswordMsg string
}

// view

// registerView 注册页面
func (server *Server) registerView(w http.ResponseWriter, r *http.Request) {
	// 是否已经登录
	server.isRedirect(w, r)
	server.renderRegister(w, r, nil)
}

// loginView 登录页面
func (server *Server) loginView(w http.ResponseWriter, r *http.Request) {
	// 是否已经登录
	server.isRedirect(w, r)
	server.renderLogin(w, r, nil)
}

// data

// register 注册
func (server *Server) register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	resp, ok := viladatorRegister(email, username, password)
	if !ok {
		server.renderRegister(w, r, resp)
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
			resp.Summary = "邮箱或名称已经存在"
			server.renderRegister(w, r, resp)
			return
		}

		resp.Summary = "请求网络错误,请刷新重试"
		server.renderRegister(w, r, resp)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

// login 登录
func (server *Server) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseForm(); err != nil {
		server.renderLogin(w, r, registerPageData{Summary: "请求网络错误,请刷新重试"})
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	resp, ok := viladatorLogin(email, password)
	if !ok {
		server.renderLogin(w, r, resp)
		return
	}

	ctx := r.Context()
	user, err := server.store.GetUserByEmail(ctx, email)
	if err != nil {
		if server.store.IsNoRows(sql.ErrNoRows) {
			resp.Summary = "邮箱或密码错误"
			server.renderLogin(w, r, resp)
			return
		}

		resp.Summary = "请求网络错误,请刷新重试"
		server.renderLogin(w, r, resp)
		return
	}

	err = pwd.BcryptComparePassword(user.HashedPassword, password)
	if err != nil {
		resp.Summary = "邮箱或密码错误"
		server.renderLogin(w, r, resp)
		return
	}

	encoded, err := server.secureCookie.Encode(AuthorizeCookie, &Authorize{ID: user.ID, Name: user.Username})
	if err != nil {
		resp.Summary = "请求网络错误,请刷新重试(cookie)"
		server.renderLogin(w, r, resp)
		return
	}

	c := cookie.NewCookie(cookie.AuthorizeName, encoded, time.Now().Add(time.Duration(7200)*time.Second))
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusFound)
}

// logout 退出
func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie.DeleteCookie(w, cookie.AuthorizeName)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// method

// renderRegister 渲染注册页面
func (server *Server) renderRegister(w http.ResponseWriter, r *http.Request, data any) {
	server.renderLayout(w, r, data, "user/register.html.tmpl")
}

// renderLogin 渲染登录页面
func (server *Server) renderLogin(w http.ResponseWriter, r *http.Request, data any) {
	server.renderLayout(w, r, data, "user/login.html.tmpl")
}

// viladatorRegister 校验注册数据
func viladatorRegister(email, username, password string) (registerPageData, bool) {
	ok := true
	resp := registerPageData{
		Email:    email,
		Username: username,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		resp.EmailMsg = "请填写正确的邮箱地址"
		ok = false
	}
	if !ValidateRxUsername(username) {
		resp.UsernameMsg = "名称(6-20,字母,数字)"
		ok = false
	}
	if !ValidatePassword(password) {
		resp.PasswordMsg = "密码(8-20位)"
		ok = false
	}

	return resp, ok
}

// viladatorLogin 校验登录数据
func viladatorLogin(email, password string) (loginPageData, bool) {
	ok := true
	errs := loginPageData{
		Email:    email,
		Password: password,
	}

	if !ValidateRxEmail(email) {
		errs.EmailMsg = "请填写正确的邮箱地址"
		ok = false
	}
	if len(password) == 0 {
		errs.PasswordMsg = "请填写正确的密码"
		ok = false
	}

	return errs, ok
}
