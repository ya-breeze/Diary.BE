package webapp

import (
	"net/http"
	"time"

	"github.com/ya-breeze/diary.be/pkg/utils"
)

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

	userID, err := r.db.GetUserID(username)
	if err != nil {
		r.logger.Warn("failed to get user ID", "username", username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := r.cookies.Get(req, "session-name")
	if err != nil {
		r.logger.Warn("failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["userID"] = userID
	// Allow to use without HTTPS - for local network
	session.Options.Secure = false
	session.Options.SameSite = http.SameSiteLaxMode
	err = session.Save(req, w)
	if err != nil {
		r.logger.Warn("failed to save session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.CreateTemplateData(req, "home")
	data["UserID"] = userID

	if err := tmpl.ExecuteTemplate(w, "home.tpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (r *WebAppRouter) logoutHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c := &http.Cookie{
		Name:     "session-name",
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
