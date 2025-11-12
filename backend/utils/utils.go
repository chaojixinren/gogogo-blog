package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gogogo/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type JWTClaims struct {
	UserID   uint   `json:"sub"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生产JWT根据用户ID和用户名
func GenerateJWT(userID uint, username string) (string, error) {
	ttl := config.AppConfig.Auth.TokenTTLHours
	if ttl <= 0 {
		ttl = 72
	}

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   strconv.FormatUint(uint64(userID), 10),
			Issuer:    config.AppConfig.App.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.Auth.JWTSecret))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.Auth.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*JWTClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
