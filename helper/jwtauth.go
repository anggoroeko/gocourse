package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     uint   `json:"role"`
	jwt.RegisteredClaims
}

//you should set it on .env file then call it from there
// var sampleSecretKey = []byte("SecretYouShouldHide")
var secretKey = os.Getenv("SECRET_KEY")
var jwtKey = []byte(secretKey)

func GenerateSecretKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil
}

// AuthMiddleware is a simple middleware to check if the request has a valid token.
func ValidateJWT() gin.HandlerFunc {
	//you should set it on .env file then call it from there
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		// Parse the token and verify its signature
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check that the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key
			return jwtKey, nil
		})

		// Check for errors
		if err != nil {
			//http.StatusUnauthorized == 401
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"error":   http.StatusUnauthorized,
			})
			c.Abort()
		}
		c.Next()
	}
}

func GenerateJWT(email, username string, role uint) (tokenString string, err error) {
	expTime := time.Now().Add(24 * time.Hour)
	// expTime := time.Now().Add(1 * time.Minute)

	claims := &JWTClaim{
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	return
}
