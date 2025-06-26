package utils

import (
	"backend/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretadmin = func() []byte {
	secret := config.GetJWTSecret("admin")
	if secret == "" {
		panic("JWT_SECRET not set in environment variables")
	}
	return []byte(secret)
}()

func GetAdminJWT(userID, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretadmin)
}

func ParseAdminJWT(tokenString string) (jwt.MapClaims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretadmin, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token validation failed: %w", jwt.ErrTokenInvalidClaims)
	}

	// Optional: Manual expiration check for extra safety
	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired: %w", jwt.ErrTokenExpired)
	}

	return claims, nil
}
