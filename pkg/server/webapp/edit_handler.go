package webapp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"github.com/ya-breeze/diary.be/pkg/utils"
)

func (r *WebAppRouter) editHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "edit")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	date := req.URL.Query().Get("date")
	if date == "" {
		date = utils.GetCurrentDate()
	}
	item, err := r.db.GetItem(userID, date)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			r.logger.Error("Failed to get item", "error", err, "date", date, "userID", userID)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item = &models.Item{
			Date:  date,
			Title: "",
			Body:  "",
		}
	}
	data["item"] = item
	data["assets"] = utils.GetAssetsFromMarkdown(item.Body)

	templateName := "edit.tpl"
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		r.logger.Warn("failed to execute template", "error", err, "template", templateName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *WebAppRouter) saveHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "edit")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	date := req.FormValue("date")
	if date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}
	item := &models.Item{
		UserID: userID,
		Date:   date,
		Title:  req.FormValue("title"),
		Body:   req.FormValue("body"),
		Tags:   strings.Split(req.FormValue("tags"), ","),
	}

	if err := r.db.PutItem(userID, item); err != nil {
		r.logger.Error("Failed to save item", "error", err, "item", item)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/?date="+item.Date, http.StatusSeeOther)
}
