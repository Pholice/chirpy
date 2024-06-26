package main

import (
	"encoding/json"
	"net/http"
)

type RequestUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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
	}
	respondWithJSON(w, http.StatusCreated, user)
}
