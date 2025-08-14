package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/diary.be/pkg/auth"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/api"
	"github.com/ya-breeze/diary.be/pkg/server/webapp"
)

func Server(logger *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("failed to open storage: %w", err)
	}

	_, finishChan, err := Serve(ctx, logger, storage, cfg)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	// Wait for an interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	logger.Info("Received signal. Shutting down server...")

	// Stop the server
	cancel()
	<-finishChan
	return nil
}

func createControllers(logger *slog.Logger, cfg *config.Config, db database.Storage, authService goserver.AuthAPIService) goserver.CustomControllers {
	return goserver.CustomControllers{
		AuthAPIService:   authService,
		UserAPIService:   api.NewUserAPIService(logger, db),
		AssetsAPIService: api.NewAssetsAPIService(logger, cfg),
	}
}

func Serve(
	ctx context.Context, logger *slog.Logger,
	storage database.Storage, cfg *config.Config,
) (net.Addr, chan int, error) {
	commit := func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}()
	logger.Info("Built from git commit: " + commit)

	if cfg.JWTSecret == "" {
		logger.Warn("JWT secret is not set. Creating random secret...")
		cfg.JWTSecret = auth.GenerateRandomString(32)
	}

	logger.Info("Starting GeekBudget server...")

	if cfg.Users != "" {
		logger.Info("Creating users...")
		users := strings.SplitSeq(cfg.Users, ",")
		for user := range users {
			tokens := strings.Split(user, ":")
			if len(tokens) != 2 {
				return nil, nil, fmt.Errorf("invalid user format: %s", user)
			}

			if err := upsertUser(storage, tokens[0], tokens[1], logger); err != nil {
				return nil, nil, fmt.Errorf("failed to update user %q: %w", tokens[0], err)
			}
		}
	} else {
		logger.Info("No users defined in configuration")
	}

	// Create shared auth service
	authService := api.NewAuthAPIService(logger, storage, cfg)

	return goserver.Serve(ctx, logger, cfg,
		createControllers(logger, cfg, storage, authService),
		[]goserver.Router{webapp.NewWebAppRouter(commit, logger, cfg, storage, authService)},
		createMiddlewares(logger, cfg)...)
}

func upsertUser(storage database.Storage, username, hashedPassword string, logger *slog.Logger) error {
	userID, err := storage.GetUserID(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("failed to reading user from DB: %w", err)
	}
	var user *models.User
	if !errors.Is(err, database.ErrNotFound) {
		logger.Info(fmt.Sprintf("Updating password for user %q", username))

		user, err = storage.GetUser(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		user.HashedPassword = hashedPassword
		if err = storage.PutUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Info(fmt.Sprintf("Creating user %q", username))
		user, err = storage.CreateUser(username, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		logger.Info(fmt.Sprintf("User %q created with ID %s", username, user.ID))
	}

	return nil
}

func createMiddlewares(logger *slog.Logger, cfg *config.Config) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		AuthMiddleware(logger, cfg),
	}
}
