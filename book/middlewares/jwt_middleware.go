// middlewares/jwt_middleware.go
package middlewares

import (
	"book/jwt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a middleware function for JWT token validation
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Token Received:", tokenString) // Add this line for debugging

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Trim any leading or trailing spaces from the token string
		// tokenString = strings.TrimSpace(tokenString)

		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
