package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ResponseUser struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	Refresh string `json:"refresh_token"`
	Red     bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestUser
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request")
		return // Ensure we return after responding
	}

	user, err := cfg.DB.GetUserEmail(reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user")
		return // Ensure we return after responding
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(reqBody.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		return // Ensure we return after responding
	}

	token, err := cfg.createJWT(user.ID, reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token")
		return // Ensure we return after responding
	}

	response := ResponseUser{
		ID:      user.ID,
		Email:   user.Email,
		Token:   token,
		Refresh: user.RefreshToken,
		Red:     user.Red,
	}
	respondWithJSON(w, http.StatusOK, response) // Correctly send the final response
}
