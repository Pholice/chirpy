package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createJWT(id int, user RequestUser) (string, error) {
	expiresInSecs := 3600
	if user.Expire > 0 && user.Expire < 3600 {
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

func (cfg *apiConfig) verifyJWT(tokenString string) (int, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		},
	)
	if err != nil {
		return 0, errors.New("COULD NOT PARSE TOKEN")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, errors.New("COULD NOT GET SUBJECT")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return 0, errors.New("COULD NOT GET ISSUER")
	}
	if issuer != string("chirpy") {
		return 0, errors.New("INVALID ISSUER")
	}
	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		return 0, errors.New("COULD NOT CONVERT ID TO INT")
	}
	return userIDInt, nil
}

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:]
		return tokenString, nil
	}
	return tokenString, errors.New("COULD NOT GET BEARER TOKEN")
}
