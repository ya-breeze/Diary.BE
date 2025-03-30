package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/server/common"
)

func AuthMiddleware(logger *slog.Logger, cfg *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Skip authorization for the root endpoint
			if req.URL.Path == "/" || strings.HasPrefix(req.URL.Path, "/web/") {
				next.ServeHTTP(writer, req)
				return
			}

			// Skip authorization for the authorize endpoint - there is no way to do it with
			// go-server openapi templates now :(
			if req.URL.Path == "/v1/authorize" {
				next.ServeHTTP(writer, req)
				return
			}

			checkToken(logger, cfg.JWTSecret, next, writer, req)
		})
	}
}

func checkToken(
	logger *slog.Logger, jwtSecret string, next http.Handler,
	writer http.ResponseWriter, req *http.Request,
) {
	// Authorization logic
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		http.Error(writer, "Invalid authorization header", http.StatusUnauthorized)
		return
	}
	bearerToken := authHeaderParts[1]

	// Parse the token
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		signingKey := []byte(jwtSecret)
		return signingKey, nil
	})
	if err != nil {
		logger.With("err", err).Warn("Invalid token")
		http.Error(writer, "Invalid token", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := claims.GetSubject()
		if err != nil {
			logger.With("err", err).Warn("Invalid subject")
			http.Error(writer, "Invalid subject", http.StatusUnauthorized)
			return
		}
		logger.With("userID", userID).Info("Authorized user")

		req = req.WithContext(context.WithValue(req.Context(), common.UserIDKey, userID))
		next.ServeHTTP(writer, req)
		return
	}

	http.Error(writer, "Unauthorized", http.StatusUnauthorized)
}
