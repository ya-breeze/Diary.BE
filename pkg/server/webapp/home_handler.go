package webapp

import (
	"net/http"

	"github.com/ya-breeze/diary.be/pkg/utils"
)

//nolint:funlen,cyclop,gocognit
func (r *WebAppRouter) homeHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "home")

	session, err := r.cookies.Get(req, "session-name")
	if err != nil {
		r.logger.Error("Failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

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
	}

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
