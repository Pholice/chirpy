package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	num, err := strconv.Atoi(chirpID)
	if err != nil {
		log.Fatal(err)
	}
	chirp, err := cfg.DB.GetChirp(num)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
