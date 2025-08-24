package api

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/common"
)

type SyncAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewSyncAPIService(logger *slog.Logger, db database.Storage) goserver.SyncAPIService {
	return &SyncAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

// GetChanges - get changes for synchronization
func (s *SyncAPIServiceImpl) GetChanges(
	ctx context.Context,
	since int32,
	limit int32,
) (goserver.ImplResponse, error) {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("User ID not found in context")
		return goserver.Response(401, nil), nil
	}

	s.logger.Info("Getting changes for sync", "userID", userID, "since", since, "limit", limit)

	// Validate parameters
	if limit <= 0 || limit > 1000 {
		limit = 100 // default limit
	}

	// Get changes from database
	sinceUint := uint(since)
	if since < 0 {
		sinceUint = 0
	}
	changes, err := s.db.GetChangesSince(userID, sinceUint, int(limit))
	if err != nil {
		s.logger.Error("Failed to get changes", "error", err, "userID", userID, "since", since)
		return goserver.Response(500, nil), nil
	}

	// Convert database changes to API response format
	responseChanges := make([]goserver.SyncChangeResponse, len(changes))
	for i, change := range changes {
		responseChanges[i] = change.ToSyncResponse()
	}

	// Determine if there are more changes available
	hasMore := len(changes) == int(limit)
	var nextID int32
	if hasMore && len(changes) > 0 {
		lastID := changes[len(changes)-1].ID
		if lastID <= uint(^uint32(0)>>1) { // Check if it fits in int32 (max positive value)
			nextID = int32(lastID) // #nosec G115 - checked above
		}
	}

	// Create the sync response
	response := goserver.SyncResponse{
		Changes: responseChanges,
		HasMore: hasMore,
		NextId:  nextID,
	}

	return goserver.Response(200, response), nil
}
