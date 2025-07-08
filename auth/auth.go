package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"os"
)

// ValidateJWT validates a JWT token and returns the user ID
func ValidateJWT(tokenString string, isRefreshToken bool) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Logger.Error("Unexpected signing method: ", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to parse JWT")
		return "", err
	}

	if !token.Valid {
		logger.Logger.Error("Invalid JWT token")
		return "", jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Logger.Error("Invalid JWT claims")
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		logger.Logger.Error("Missing or invalid user_id in JWT claims")
		return "", jwt.ErrTokenInvalidId
	}

	if isRefreshToken {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			logger.Logger.Error("Invalid token type for refresh token")
			return "", jwt.ErrTokenInvalidClaims
		}
	}

	return userID, nil
}
