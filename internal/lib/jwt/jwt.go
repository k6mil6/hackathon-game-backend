package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(id int, username string, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(duration).Unix()

	return token.SignedString([]byte(secret))
}

func GetID(jwtToken string, secret string) (int, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("token is invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("ID claim is not a number")
	}

	return int(idFloat), nil
}
