package middleware

import (
	"backend/pkg/utils"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Periksa apakah token ada di blacklist
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM blacklisted_tokens WHERE token = $1 LIMIT 1)", tokenString).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			c.Abort()
			return
		}

		if exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}

		// Verifikasi token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
