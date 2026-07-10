package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"

	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// Claims extends standard jwt.Claims
type Claims struct {
	jwt.RegisteredClaims
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Type     TokenType `json:"type"`
}
