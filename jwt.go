package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createJWT(id int, user RequestUser) (string, error) {
	expiresInSecs := 86400
	if user.Expire > 0 && user.Expire < 86400 {
		expiresInSecs = user.Expire // 24 hours in seconds
	}

	expire := jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(time.Second * time.Duration(expiresInSecs))))
	issuedTime := jwt.NewNumericDate(time.Now().UTC())
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedTime,
		ExpiresAt: expire,
		Subject:   fmt.Sprintf("%d", id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", errors.New("COULD NOT GEN TOKEN")
	}
	return ss, nil
}
