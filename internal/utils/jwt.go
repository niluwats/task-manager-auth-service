package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtClaims struct {
	ID    int
	Email string
	jwt.StandardClaims
}

func GenerateJWT(email string, ID uint) (string, error) {
	expiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	claims := &jwtClaims{
		ID:    int(ID),
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiration)).Unix(),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(signedToken string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
