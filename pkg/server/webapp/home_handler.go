package webapp

import (
	"context"
	"errors"
	"html/template"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/common"
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

	// Get items data and populate template
	if err := r.populateItemsData(data, userID, date, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	r.renderTemplate(w, tmpl, data)
}

// populateItemsData fetches items data and populates the template data
func (r *WebAppRouter) populateItemsData(data map[string]any, userID, date string, req *http.Request) error {
	// Create context with user ID for the items service
	ctx := context.WithValue(req.Context(), common.UserIDKey, userID)

	// Use the items service to get the item
	response, err := r.itemsService.GetItems(ctx, date)
	if err != nil {
		r.logger.Error("Failed to get items from service", "error", err, "date", date, "userID", userID)
		return err
	}

	if response.Code != 200 {
		r.logger.Error("Items service returned non-200 status", "code", response.Code, "date", date, "userID", userID)
		return errors.New("failed to get items")
	}

	itemsResponse, ok := response.Body.(goserver.ItemsResponse)
	if !ok {
		r.logger.Error("Failed to cast response body to ItemsResponse")
		return errors.New("internal server error")
	}

	// Convert the service response to template data
	data["item"] = map[string]any{
		"Date":  itemsResponse.Date,
		"Title": itemsResponse.Title,
		"Body":  itemsResponse.Body,
		"Tags":  itemsResponse.Tags,
	}

	body := markdown.ToHTML([]byte(itemsResponse.Body), nil, utils.NewImagePrefixRenderer("/web/assets/"))
	//nolint:gosec // this is safe
	data["body"] = template.HTML(string(body))

	// Add navigation dates from the service response
	if itemsResponse.PreviousDate != nil {
		data["previousDate"] = *itemsResponse.PreviousDate
	}
	if itemsResponse.NextDate != nil {
		data["nextDate"] = *itemsResponse.NextDate
	}

	return nil
}

// renderTemplate renders the template with the provided data
func (r *WebAppRouter) renderTemplate(w http.ResponseWriter, tmpl *template.Template, data map[string]any) {
	// if utils.IsMobile(req.Header.Get("User-Agent")) {
	// data["Template"] = "home_mobile.tpl"
	// } else {
	data["Template"] = "home.tpl"
	// }

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
