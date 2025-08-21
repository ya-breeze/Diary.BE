package database

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"gorm.io/gorm"
)

//go:generate go tool github.com/golang/mock/mockgen -destination=mocks/mock_storage.go -package=mocks github.com/ya-breeze/diary.be/pkg/database Storage

const StorageError = "storage error: %w"

var ErrNotFound = errors.New("not found")

// SearchParams defines parameters for searching diary items
type SearchParams struct {
	// SearchText filters items by title and body content (case-insensitive)
	SearchText string
	// Tags filters items that contain any of the specified tags
	Tags []string
	// Date filters items by specific date (optional, for backward compatibility)
	Date string
}

//nolint:interfacebloat // keep a single storage interface for simplicity
type Storage interface {
	Open() error
	Close() error

	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
	PutUser(user *models.User) error

	GetItem(userID, itemID string) (*models.Item, error)
	GetItems(userID string, searchParams SearchParams) ([]*models.Item, int, error)
	PutItem(userID string, item *models.Item) error

	GetPreviousDate(userID, date string) (string, error)
	GetNextDate(userID, date string) (string, error)
}

type storage struct {
	log *slog.Logger
	cfg *config.Config
	db  *gorm.DB
}

func NewStorage(logger *slog.Logger, cfg *config.Config) Storage {
	return &storage{log: logger, db: nil, cfg: cfg}
}

func (s *storage) Open() error {
	s.log.Info("Opening database", "path", s.cfg.DBPath)
	var err error
	s.db, err = openSqlite(s.log, s.cfg.DBPath, s.cfg.Verbose)
	if err != nil {
		s.log.Error("failed to connect database", "error", err)
		panic("failed to connect database")
	}
	if err := autoMigrateModels(s.db); err != nil {
		s.log.Error("failed to migrate database", "error", err)
		panic("failed to migrate database")
	}

	return nil
}

func (s *storage) Close() error {
	// return s.db.Close()
	return nil
}

func (s *storage) CreateUser(username, hashedPassword string) (*models.User, error) {
	_, err := s.GetUserID(username)
	if err == nil {
		s.log.Error("user already exists", "username", username)
		return nil, fmt.Errorf("user %q already exists", username)
	}

	user := models.User{
		ID:             uuid.New(),
		Login:          username,
		HashedPassword: hashedPassword,
		StartDate:      time.Now(),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	return &user, nil
}

func (s *storage) GetUser(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf(StorageError, err)
	}

	return &user, nil
}

func (s *storage) PutUser(user *models.User) error {
	existingUserID, err := s.GetUserID(user.Login)
	if err != nil {
		s.log.Error("failed to get user ID", "error", err, "user", user.Login)
		return fmt.Errorf("failed to get user ID: %w", err)
	}
	if existingUserID != user.ID.String() {
		s.log.Error("user ID mismatch", "expected", user.ID.String(), "actual", existingUserID)
		return fmt.Errorf("user ID mismatch: expected %s, actual %s", user.ID.String(), existingUserID)
	}

	// Update the user in the database
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetUserID(username string) (string, error) {
	var user models.User
	if err := s.db.Where("login = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf(StorageError, err)
	}

	return user.ID.String(), nil
}

// #region Item

func (s *storage) GetItem(userID, date string) (*models.Item, error) {
	var item models.Item
	if err := s.db.Where("date = ? and user_id = ?", date, userID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf(StorageError, err)
	}

	return &item, nil
}

func (s *storage) GetItems(userID string, searchParams SearchParams) ([]*models.Item, int, error) {
	var items []*models.Item
	query := s.db.Where("user_id = ?", userID)

	// Apply date filter if specified (for backward compatibility)
	if searchParams.Date != "" {
		query = query.Where("date = ?", searchParams.Date)
	}

	// Apply text search filter if specified
	if searchParams.SearchText != "" {
		searchPattern := "%" + searchParams.SearchText + "%"
		query = query.Where("title LIKE ? OR body LIKE ?", searchPattern, searchPattern)
	}

	// Apply tag filters if specified
	if len(searchParams.Tags) > 0 {
		// For JSON tag filtering, we need to check if any of the specified tags exist in the JSON array
		tagConditions := make([]string, len(searchParams.Tags))
		tagArgs := make([]any, len(searchParams.Tags))
		for i, tag := range searchParams.Tags {
			tagConditions[i] = "JSON_EXTRACT(tags, '$') LIKE ?"
			tagArgs[i] = "%\"" + tag + "\"%"
		}
		tagQuery := strings.Join(tagConditions, " OR ")
		query = query.Where(tagQuery, tagArgs...)
	}

	// Get total count for pagination
	var totalCount int64
	if err := query.Model(&models.Item{}).Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf(StorageError, err)
	}

	// Execute the query to get items, ordered by date descending
	if err := query.Order("date DESC").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf(StorageError, err)
	}

	return items, int(totalCount), nil
}

func (s *storage) PutItem(userID string, item *models.Item) error {
	item.UserID = userID
	if err := s.db.Save(item).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

// #endregion Item

// #region Dates

func (s *storage) GetPreviousDate(userID, date string) (string, error) {
	var item models.Item
	if err := s.db.Where("user_id = ? and date < ?", userID, date).Order("date desc").First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf(StorageError, err)
	}

	return item.Date, nil
}

func (s *storage) GetNextDate(userID, date string) (string, error) {
	var item models.Item
	if err := s.db.Where("user_id = ? and date > ?", userID, date).Order("date asc").First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf(StorageError, err)
	}

	return item.Date, nil
}

// #endregion Dates
