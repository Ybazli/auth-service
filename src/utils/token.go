package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ybazli/auth-service/src/models"
	"os"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
}

func GenerateTokenPair(user models.User) (TokenPair, error) {
	sessionID := uuid.NewString()
	accessTokenClaim := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)

	accessSecret := os.Getenv("JWT_SECRET")
	accessTokenStr, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken := uuid.NewString()

	return TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}, nil
}
