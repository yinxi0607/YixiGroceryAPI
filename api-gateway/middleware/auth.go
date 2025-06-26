package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/auth/register" || c.Request.URL.Path == "/api/auth/login" {
			c.Next()
			return
		}

		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(401, gin.H{"code": 401, "message": "Missing token"})
			c.Abort()
			return
		}

		if len(tokenStr) > 7 && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_jwt_secret"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid claims"})
			c.Abort()
			return
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid user_id"})
			c.Abort()
			return
		}
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
