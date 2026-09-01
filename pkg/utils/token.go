package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// Claims represents standard JWT claims with application public identity
type Claims struct {
	jwt.RegisteredClaims
	PublicID string   `json:"public_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles,omitempty"`
}

func NewClaims(
	PublicID string, Email string, Roles []string, expiryDuration time.Duration,
) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-svc",
		},
		PublicID: PublicID,
		Email:    Email,
		Roles:    Roles,
	}
}

// ParseRSAPrivateKey parses an RSA private key from PEM formatted bytes.
func ParseRSAPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)
	if err != nil {
		return nil, fmt.Errorf("parse rsa private key: %w", err)
	}
	return key, nil
}

// ParseRSAPublicKey parses an RSA public key from PEM formatted bytes.
func ParseRSAPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
	if err != nil {
		return nil, fmt.Errorf("parse rsa public key: %w", err)
	}
	return key, nil
}

// SignRS256 signs claims using an RSA private key with the RS256 algorithm.
func SignRS256(claims Claims, privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", errors.New("rsa private key is required for signing")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign rs256 token: %w", err)
	}
	return tokenString, nil
}

// VerifyRS256 verifies a JWT string using an RSA public key and returns the parsed Claims.
func VerifyRS256(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	if publicKey == nil {
		return nil, errors.New("rsa public key is required for verification")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	},
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer("auth-svc"),
	)

	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GenerateRandomToken generates a cryptographically secure hex-encoded random token.
func GenerateRandomToken(bytesLen int) (string, error) {
	if bytesLen <= 0 {
		bytesLen = 32
	}
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
