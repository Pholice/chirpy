package main

import "net/http"

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := getBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not get bearer token")
	}
	err = cfg.DB.RevokeRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not remove token")
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
