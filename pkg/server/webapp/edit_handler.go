package webapp

import (
	"net/http"

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

	session, err := r.cookies.Get(req, "session-name")
	if err != nil {
		r.logger.Error("Failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["userID"].(string)
	if !ok {
		if err := tmpl.ExecuteTemplate(w, "login.tpl", data); err != nil {
			r.logger.Warn("failed to execute login template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	data["UserID"] = userID

	itemID := req.URL.Query().Get("itemID")
	data["itemID"] = itemID

	if itemID != "" {
		item, err := r.db.GetItem(userID, itemID)
		if err != nil {
			r.logger.Error("Failed to get item", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["item"] = item
	} else {
		data["item"] = &models.Item{
			Title: "",
			Text:  "",
		}
	}

	// dateFrom, dateTo, err := getTimeRange(req, utils.GranularityMonth)
	// if err != nil {
	// 	r.logger.Error("Failed to get time range", "error", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// dateFrom = dateFrom.AddDate(0, -12, 0)

	// data["From"] = dateFrom
	// data["To"] = dateTo
	// data["Current"] = dateFrom.Unix()
	// data["Last"] = time.Date(
	// 	dateFrom.Year(), 1, 1, 0, 0, 0, 0, dateFrom.Location(),
	// ).AddDate(-1, 0, 0).Unix()
	// data["Next"] = dateTo.Unix()

	templateName := "edit.tpl"
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		r.logger.Warn("failed to execute template", "error", err, "template", templateName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
