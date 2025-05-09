package jwt

import (
	"time"

	"github.com/takuchi17/term-keeper/app/models"
	"github.com/takuchi17/term-keeper/configs"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(configs.Config.JWTSecret)

func GenerateToken(userID models.UserId, userName models.UserName) (string, error) {
	claims := jwt.MapClaims{
		"userid":   userID,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
