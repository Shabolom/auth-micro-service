package utils

import (
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
		logger.Info("Token is invalid", zap.String("token", strToken))
		return Claims{}, err
	}

	return claims, nil
}
