package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not get token")
		return
	}
	userIDInt, err := cfg.verifyJWT(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user")
		return
	}

	var reqBody RequestUser
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid request body")
		return
	}
	user, err := cfg.DB.UpdateUser(userIDInt, reqBody.Email, reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	response := payload{
		ID:    user.ID,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, response)
}
