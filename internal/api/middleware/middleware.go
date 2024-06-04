package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"pir-serverSide/config"
	"pir-serverSide/internal/repository/models"
	"time"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": "вы не авторизованны",
		})
		return
	}
	if err != nil {
		c.Status(http.StatusUnauthorized)
	}
	env, _ := config.LoadConfig()
	secret := env.Secret

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Status(http.StatusUnauthorized)
		}
		var user models.User
		config.DB.First(&user, claims["sub"]).Preload("Applications")

		if user.Id == 0 {
			c.Status(http.StatusUnauthorized)
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.Status(http.StatusUnauthorized)
	}

}
