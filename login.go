package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestUser
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode")
	}
	user, err := cfg.DB.GetUser(reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user")
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(reqBody.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
	}
	respondWithJSON(w, http.StatusOK, user)
}
