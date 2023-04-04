package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"

	helpers "github.com/crazi-coder/report-service/core/utils/helpers"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(helpers.GetEnv("JWT_SECRET", "lqEjTETjq0vETXloAKJcFKlGSan9OgPVaX3LYBnwJPNhNGFPEfWUjadpmkyyG1sG"))

// JWTClaims custom declaration structure and embedded JWT StandardClaims
// jwt package comes with jwt Standardclaims contains only official fields
// We need to record an additional user information field here, so we need to customize the structure
// If you want to save more information, you can add it to this structure
type JWTClaims struct {
	UserID   string   `json:"id"`
	UserRole []string `json:"roles"`
	Schema   string   `json:"uk"`
	jwt.RegisteredClaims
}

// parseToken parsing JWT
func ParseToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid { // Verification token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func CreateToken(ctx context.Context, UserID string, UserGroup []string, Schema string, domain string) (string, error) {
	st := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "Infilect Pvt Ltd",
		Subject:   "AuthToken",
		ID:        "1298JUW",
		Audience:  []string{domain},
	}
	claims := JWTClaims{
		UserID:           UserID,
		UserRole:         UserGroup,
		Schema:           Schema,
		RegisteredClaims: st,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
