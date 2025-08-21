package webapp

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// List of allowed file extensions
//
//nolint:gochecknoglobals
var allowedExtensions = []string{
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".mp4", ".mov", ".avi", ".wmv", ".flv", ".mkv",
}

func (r *WebAppRouter) uploadHandler(w http.ResponseWriter, req *http.Request) {
	userID, code, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err = req.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	asset, header, err := req.FormFile("asset")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer asset.Close()

	// Save the file to the server
	userAssetPath := filepath.Join(r.cfg.AssetPath, userID)
	if err = os.MkdirAll(userAssetPath, 0o755); err != nil {
		r.logger.Error("Failed to create directory", "error", err, "path", userAssetPath)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check file by extension from header - only allow images and videos
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !slices.Contains(allowedExtensions, ext) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Create the uploaded file on the server
	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filePath := filepath.Join(userAssetPath, filename)
	r.logger.Info("Saving file", "path", filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		r.logger.Error("Failed to create file", "error", err, "path", filePath)
		http.Error(w, "Could not save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's data to the server file
	if _, err := io.Copy(dst, asset); err != nil {
		r.logger.Error("Failed to copy file", "error", err)
		http.Error(w, "Could not write the file", http.StatusInternalServerError)
		return
	}

	// Respond with the path or URL to the uploaded file
	fmt.Fprint(w, filename)
}
