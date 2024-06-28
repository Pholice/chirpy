package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RequestChirp struct {
	Body string `json:"body"`
}

func filter(words []string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	var copy string = ""
	for _, word := range words {
		isBad := false
		for _, bad := range badWords {
			if strings.ToLower(word) == bad {
				copy += "**** "
				isBad = true
				break
			}
		}
		if !isBad {
			copy += word + " "
		}
	}
	return strings.TrimSpace(copy)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestChirp
	tokenString, err := getBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad token format")
	}
	userIDInt, err := cfg.verifyJWT(tokenString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not verify JWT token")
		return
	}
	user, err := cfg.DB.GetUserID(userIDInt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find user")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(reqBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "error: Chirp is too long")
		return
	}

	clean := filter(strings.Split(reqBody.Body, " "))
	chirp, err := cfg.DB.CreateChirp(clean, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}
