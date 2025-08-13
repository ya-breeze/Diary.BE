package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/common"
)

type AssetsAPIServiceImpl struct {
	logger *slog.Logger
	cfg    *config.Config
}

func NewAssetsAPIService(logger *slog.Logger, cfg *config.Config) goserver.AssetsAPIService {
	return &AssetsAPIServiceImpl{
		logger: logger,
		cfg:    cfg,
	}
}

// GetAsset - return asset by path
func (s *AssetsAPIServiceImpl) GetAsset(ctx context.Context, path string) (goserver.ImplResponse, error) {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("Failed to get user ID from context")
		return goserver.Response(http.StatusUnauthorized, nil), nil
	}

	// Validate path to prevent directory traversal attacks
	if strings.Contains(path, "..") {
		s.logger.Warn("Invalid asset path requested (contains ..)", "path", path, "userID", userID)
		return goserver.Response(http.StatusBadRequest, nil), nil
	}

	// Clean the path to normalize separators and remove any redundant elements
	cleanPath := filepath.Clean(path)

	// Ensure the clean path doesn't start with / or \ (absolute path)
	if filepath.IsAbs(cleanPath) {
		s.logger.Warn("Invalid asset path requested (absolute path)", "path", path, "userID", userID)
		return goserver.Response(http.StatusBadRequest, nil), nil
	}

	// Construct the full path to the user's asset
	userAssetBasePath := filepath.Join(s.cfg.AssetPath, userID)
	userAssetPath := filepath.Join(userAssetBasePath, cleanPath)

	// Ensure the resolved path is still within the user's asset directory
	absBasePath, err := filepath.Abs(userAssetBasePath)
	if err != nil {
		s.logger.Error("Failed to get absolute base path", "error", err, "basePath", userAssetBasePath)
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	absAssetPath, err := filepath.Abs(userAssetPath)
	if err != nil {
		s.logger.Error("Failed to get absolute asset path", "error", err, "assetPath", userAssetPath)
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	// Check if the resolved path is within the user's directory
	if !strings.HasPrefix(absAssetPath, absBasePath+string(filepath.Separator)) && absAssetPath != absBasePath {
		s.logger.Warn("Asset path outside user directory", "path", path, "resolvedPath", absAssetPath, "userID", userID)
		return goserver.Response(http.StatusBadRequest, nil), nil
	}

	s.logger.Info("Serving asset", "path", userAssetPath, "userID", userID)

	// Check if file exists and is accessible
	fileInfo, err := os.Stat(userAssetPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Debug("Asset not found", "path", userAssetPath, "userID", userID)
			return goserver.Response(http.StatusNotFound, nil), nil
		}
		s.logger.Error("Failed to stat asset file", "error", err, "path", userAssetPath, "userID", userID)
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	// Ensure it's a file, not a directory
	if fileInfo.IsDir() {
		s.logger.Warn("Requested path is a directory", "path", userAssetPath, "userID", userID)
		return goserver.Response(http.StatusBadRequest, nil), nil
	}

	// Open the file
	file, err := os.Open(userAssetPath)
	if err != nil {
		s.logger.Error("Failed to open asset file", "error", err, "path", userAssetPath, "userID", userID)
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	// Return the file - the framework will handle closing it and setting appropriate headers
	return goserver.Response(http.StatusOK, file), nil
}
