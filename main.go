package main

import (
	"log"
	"net/http"

	"github.com/Pholice/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func main() {
	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	serveMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serveMux.Handle("GET /api/healthz", http.HandlerFunc(apiCfg.healthz))
	serveMux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.metrics))
	serveMux.Handle("GET /api/reset", http.HandlerFunc(apiCfg.reset))
	serveMux.Handle("POST /api/chirps", http.HandlerFunc(apiCfg.createChirp))
	serveMux.Handle("POST /api/users", http.HandlerFunc(apiCfg.createUser))
	serveMux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(apiCfg.getChirp))

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
