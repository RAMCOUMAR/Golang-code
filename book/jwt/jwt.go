// jwt/jwt.go
package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with the provided user information
func GenerateToken(userID int, username string) (string, error) {
	// Load the secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not set")
	}

	// Using github.com/dgrijalva/jwt-go library
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Set the expiration time as needed
		},
	})

	// Convert the string secret key to a byte slice
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken verifies the JWT token and returns the claims
func VerifyToken(tokenString string) (*Claims, error) {
	// Load the secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY not set")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
