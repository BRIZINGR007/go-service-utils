package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Email  string
	UserID string
}

var secretKey = os.Getenv("JWT_SECRET_KEY")

func GenerateToken(email string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (*TokenClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected Signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email claim missing or invalid")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return nil, errors.New("userId claim missing or invalid")
	}
	return &TokenClaims{
		Email:  email,
		UserID: userId,
	}, nil

}
