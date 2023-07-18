package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTWrapper struct {
	SecretKey  string
	Issuer     string
	Expiration int
}

type jwtClaims struct {
	ID    int
	Email string
	jwt.StandardClaims
}

func (wrapper *JWTWrapper) GenerateJWT(email string, ID int) (string, error) {
	claims := &jwtClaims{
		ID:    ID,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(wrapper.Expiration)).Unix(),
			Issuer:    wrapper.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(wrapper.SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
