package database

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"gorm.io/gorm"
)

//go:generate go tool github.com/golang/mock/mockgen -destination=mocks/mock_storage.go -package=mocks github.com/ya-breeze/diary.be/pkg/database Storage

const StorageError = "storage error: %w"

var ErrNotFound = errors.New("not found")

type Storage interface {
	Open() error
	Close() error

	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
	PutUser(user *models.User) error

	GetItem(userID, itemID string) (*models.Item, error)
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
