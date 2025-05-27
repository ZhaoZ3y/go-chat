package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte("s7F3qP9kLm2RtV5XwY8zJ6nHcD1Gb4NvKxQeA0TyIuCf7oM5lWpZs8EdBhO")
var refreshSecret = []byte("MIIEowIBAAKCAQEAxKf7l4J7Z8q9X2Wn1P3mCL5sYkGp2dHj6vQyT0zN1w")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成短期访问 token
func GenerateAccessToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-chat-IM",
			Subject:   "access-token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken 生成长期 refresh token
func GenerateRefreshToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-chat-IM",
			Subject:   "refresh-token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ParseAccessToken 解析访问 token
func ParseAccessToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, accessSecret)
}

// ParseRefreshToken 解析 refresh token
func ParseRefreshToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, refreshSecret)
}

// 内部通用解析函数
func parseToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
