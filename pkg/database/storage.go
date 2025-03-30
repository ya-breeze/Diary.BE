package database

import (
	"errors"
	"log/slog"

	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"gorm.io/gorm"
)

//go:generate go tool github.com/golang/mock/mockgen -destination=mocks/mock_storage.go -package=mocks github.com/ya-breeze/diary.be/pkg/database Storage

const StorageError = "storage error: %w"

var ErrNotFound = errors.New("not found")

//nolint:interfacebloat
type Storage interface {
	Open() error
	Close() error

	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
	PutUser(user *models.User) error
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
