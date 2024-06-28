package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpIDInt, err := strconv.Atoi(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not convert string")
		return
	}
	chirp, err := cfg.DB.GetChirp(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not get chirp")
		return
	}

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

	user, err := cfg.DB.GetUserID(userIDInt)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not convert string")
		return
	}
	if user.ID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "Not author of chirp")
		return
	}
	err = cfg.DB.DeleteChirp(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not delete chirp")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
