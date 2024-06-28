package main

import (
	"net/http"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find refresh token")
		return
	}
	user, err := cfg.DB.UserForRefreshToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user with token")
		return
	}
	reqUser := RequestUser{
		Email:    user.Email,
		Password: string(user.Password),
		Expire:   user.ExpiresAt.Day(),
	}
	accessToken, err := cfg.createJWT(user.ID, reqUser)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}
	type payload struct {
		Token string `json:"token"`
	}
	send := payload{
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, send)
}
