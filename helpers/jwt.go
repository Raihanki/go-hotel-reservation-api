package helpers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	signingToken, err := token.SignedString([]byte(configs.ENV.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return signingToken, nil
}

func ValidateJWT(token string) (*jwt.Token, error) {
	validatedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return []byte(configs.ENV.JWT_SECRET), nil
	})
	if err != nil {
		log.Error("error validating token : ", err.Error())
		return nil, err
	}

	return validatedToken, nil
}
