package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/auth/register" || c.Request.URL.Path == "/api/auth/login" {
			c.Next()
			return
		}

		tokenStr := c.GetHeader("content")
		if tokenStr == "" {
			c.JSON(401, gin.H{"code": 401, "message": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_jwt_secret"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}
