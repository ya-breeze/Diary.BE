package webapp

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"github.com/ya-breeze/diary.be/pkg/utils"
)

func (r *WebAppRouter) homeHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "home")

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
	body := markdown.ToHTML([]byte(item.Body), nil, utils.NewImagePrefixRenderer("/web/assets/"))
	//nolint:gosec // this is safe
	data["body"] = template.HTML(string(body))

	if utils.IsMobile(req.Header.Get("User-Agent")) {
		data["Template"] = "home_mobile.tpl"
	} else {
		data["Template"] = "home.tpl"
	}

	templateName, ok := data["Template"].(string)
	if !ok {
		r.logger.Error("Failed to assert template name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		r.logger.Warn("failed to execute template", "error", err, "template", templateName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
