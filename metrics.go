package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	resp, err := os.ReadFile("metrics.html")
	if err != nil {
		fmt.Printf("Could not read html")
	}
	content := fmt.Sprintf((string(resp)), cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
