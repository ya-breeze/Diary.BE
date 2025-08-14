package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/common"
)

type ItemsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewItemsAPIService(logger *slog.Logger, db database.Storage) goserver.ItemsAPIService {
	return &ItemsAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

// GetItems - get diary items
func (s *ItemsAPIServiceImpl) GetItems(ctx context.Context, date string) (goserver.ImplResponse, error) {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("User ID not found in context")
		return goserver.Response(401, nil), nil
	}

	s.logger.Info("Getting items", "userID", userID, "date", date)

	// Get the item for the specified date
	item, err := s.db.GetItem(userID, date)

	var response goserver.ItemsResponse
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// Return empty item if not found
			response = goserver.ItemsResponse{
				Date:  date,
				Title: "",
				Body:  "",
				Tags:  []string{},
			}
		} else {
			s.logger.Error("Failed to get item", "error", err, "userID", userID, "date", date)
			return goserver.Response(500, nil), nil
		}
	} else {
		// Convert database item to API response
		response = goserver.ItemsResponse{
			Date:  item.Date,
			Title: item.Title,
			Body:  item.Body,
			Tags:  []string(item.Tags),
		}
	}

	// Add navigation dates (common for both found and not found items)
	s.addNavigationDates(&response, userID, date)

	return goserver.Response(200, response), nil
}

// addNavigationDates adds previous and next dates to the response
func (s *ItemsAPIServiceImpl) addNavigationDates(response *goserver.ItemsResponse, userID, date string) {
	if previousDate, err := s.db.GetPreviousDate(userID, date); err == nil {
		response.PreviousDate = &previousDate
	}
	if nextDate, err := s.db.GetNextDate(userID, date); err == nil {
		response.NextDate = &nextDate
	}
}
