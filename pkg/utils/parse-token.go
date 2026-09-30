package utils

import (
	"auth-micro-service/pkg/shortcut"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func ParseToken(strToken string, secret string, logger *zap.Logger) (Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(strToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не подходящий алгоритм шифрования: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return Claims{}, shortcut.ErrInvalidToken
	}

	return claims, nil
}

func ParseTokenForCheck(strToken string, secret string, logger *zap.Logger) (Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(strToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не подходящий алгоритм шифрования: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return claims, shortcut.ErrInvalidToken
	}

	return claims, nil
}
