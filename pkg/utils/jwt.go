package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your-secret-key")

func GenerateJWT(userID int, username, role string) (string, error) {
	// 生产环境使用更短的过期时间（如15分钟）
	expirationTime := time.Now().Add(15 * time.Minute)
	if os.Getenv("ENVIRONMENT") == "development" {
		expirationTime = time.Now().Add(24 * time.Hour) // 开发环境可延长
	}

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "mybot-app", // 添加签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
