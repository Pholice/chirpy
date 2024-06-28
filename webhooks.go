package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type webhook struct {
	Event string         `json:"event"`
	Data  map[string]int `json:"data"`
}

func (cfg *apiConfig) webhooks(w http.ResponseWriter, r *http.Request) {
	key, err := getAPIToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could get token")
		return
	}
	if key != cfg.Polka {
		respondWithError(w, http.StatusUnauthorized, "Incorrect apikey")
		return
	}

	var reqBody webhook
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request")
		return
	}
	if reqBody.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event")
		return
	}
	_, err = cfg.DB.UpdateUserType(reqBody.Data["user_id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not update to red")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

func getAPIToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := ""
	if strings.HasPrefix(authHeader, "ApiKey ") {
		tokenString = authHeader[7:]
		return tokenString, nil
	}
	return tokenString, errors.New("COULD NOT GET API TOKEN")
}
