package main

import (
	"encoding/json"
	"net/http"
)

type RequestUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Expire   int    `json:"expires_in_seconds"`
}
type payload struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Red   bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestUser
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	user, err := cfg.DB.CreateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	response := payload{
		ID:    user.ID,
		Email: user.Email,
		Red:   user.Red,
	}
	respondWithJSON(w, http.StatusCreated, response)
}
