package auth

import (
	"math/rand/v2"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt with a generated salt
func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// CheckPasswordHash compares the hashed password with the plain password
func CheckPasswordHash(password, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))] //nolint:gosec // it's okay to use math/rand here
	}
	return string(b)
}

func CreateJWT(userID, secret string) (string, error) {
	signingKey := []byte(secret)

	claims := &jwt.RegisteredClaims{
		Issuer:    "geekbudget",
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
