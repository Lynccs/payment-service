package middleware

import (
	"net/http"

	"github.com/Lynccs/payment-service/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := jwt.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
