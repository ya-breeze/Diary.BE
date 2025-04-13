package webapp

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (r *WebAppRouter) assetsHandler(w http.ResponseWriter, req *http.Request) {
	session, err := r.cookies.Get(req, "session-name")
	if err != nil {
		r.logger.Error("Failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	userAsset := filepath.Join(r.cfg.AssetPath, userID, strings.TrimPrefix(req.URL.Path, "/web/assets/"))
	r.logger.Info("Serving asset", "path", userAsset)
	http.ServeFile(w, req, userAsset)
}
