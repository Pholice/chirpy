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

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	asc := true
	if sort == "dec" {
		asc = false
	}
	if s != "" {
		sInt, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Could not convert to int")
			return
		}
		chirps, err := cfg.DB.GetChirpsAuthor(sInt, asc)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Could not find chirps")
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	chirps, err := cfg.DB.GetChirps(asc)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
