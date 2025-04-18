package webapp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ya-breeze/diary.be/pkg/auth"
	"github.com/ya-breeze/diary.be/pkg/utils"
)

//nolint:cyclop,funlen // refactor
func (r *WebAppRouter) loginHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := req.Form.Get("username")
	password := req.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// check that user exists in DB
	userID, err := r.db.GetUserID(username)
	if err != nil {
		r.logger.Warn("failed to get user ID", "username", username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := r.db.GetUser(userID)
	if err != nil {
		r.logger.Warn("failed to get user", "ID", userID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user == nil {
		r.logger.Warn("user not found", "ID", userID)
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// check that password is correct and create JWT token
	hashed, err := base64.StdEncoding.DecodeString(user.HashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !auth.CheckPasswordHash([]byte(password), hashed) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	token, err := auth.CreateJWT(userID, r.cfg.Issuer, r.cfg.JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// set JWT token in cookie
	session, err := r.cookies.Get(req, r.cfg.CookieName)
	if err != nil {
		r.logger.Warn("failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["token"] = token
	// Allow to use without HTTPS - for local network
	session.Options.Secure = false
	session.Options.SameSite = http.SameSiteLaxMode
	err = session.Save(req, w)
	if err != nil {
		r.logger.Warn("failed to save session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *WebAppRouter) logoutHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := &http.Cookie{
		Name:     r.cfg.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.CreateTemplateData(req, "login")

	if err := tmpl.ExecuteTemplate(w, "login.tpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *WebAppRouter) GetUserIDFromSession(req *http.Request) (string, int, error) {
	session, err := r.cookies.Get(req, r.cfg.CookieName)
	if err != nil {
		r.logger.Error("Failed to get session", "error", err)
		return "", http.StatusBadRequest, err
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		r.logger.Warn("failed to get token from session")
		return "", http.StatusBadRequest, errors.New("token not found in session")
	}

	userID, err := auth.CheckJWT(token, r.cfg.Issuer, r.cfg.JWTSecret)
	if err != nil {
		r.logger.With("err", err).Warn("Invalid token")
		return "", http.StatusUnauthorized, err
	}

	return userID, http.StatusOK, nil
}

func (r *WebAppRouter) ValidateUserID(
	tmpl *template.Template, w http.ResponseWriter, req *http.Request,
) (string, error) {
	userID, _, err := r.GetUserIDFromSession(req)
	if err != nil {
		if errTmpl := tmpl.ExecuteTemplate(w, "login.tpl", nil); errTmpl != nil {
			r.logger.Warn("failed to execute login template", "error", errTmpl)
			http.Error(w, errTmpl.Error(), http.StatusInternalServerError)
		}
		return "", fmt.Errorf("failed to get user ID from session: %w", err)
	}

	return userID, nil
}
