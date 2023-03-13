package helpers

import (
	"errors"
	"log"
	"os"
	_ "strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GenerateToken(id string, email string) string {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         id,
		"email":      email,
		"exp":        time.Now().Add(time.Minute * 5).Unix(),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if err := godotenv.Load("../api-mygram-go/config/env/.env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	signedToken, _ := parseToken.SignedString([]byte(os.Getenv("TOKEN_KEY")))

	return signedToken
}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	errResponse := errors.New("sign in to proceed")
	// headerToken := ctx.Request.Header.Get("Authorization")
	// bearer := strings.HasPrefix(headerToken, "Bearer")

	// if !bearer {
	// 	return nil, errResponse
	// }

	// stringToken := strings.Split(headerToken, " ")[1]

	// if err := godotenv.Load("../api-mygram-go/config/env/.env"); err != nil {
	// 	log.Fatal("Error loading .env file: ", err)
	// }

	// token, _ := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, errResponse
	// 	}

	// 	return []byte(os.Getenv("TOKEN_KEY")), nil
	// })

	// if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
	// 	return nil, errResponse
	// }

	// return token.Claims.(jwt.MapClaims), nil

	if err := godotenv.Load("../api-mygram-go/config/env/.env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	tokenString, err := ctx.Cookie("Authorization")

	if err != nil {
		return nil, errResponse
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}

		return []byte(os.Getenv("TOKEN_KEY")), nil
	})

	if err != nil {
		return nil, errResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errResponse
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return nil, errResponse
	}

	return claims, nil
}
