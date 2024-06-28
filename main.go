package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Pholice/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	Secret         string
	Polka          string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		Secret:         jwtSecret,
		Polka:          polkaKey,
	}

	serveMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serveMux.Handle("GET /api/healthz", http.HandlerFunc(apiCfg.healthz))
	serveMux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.metrics))
	serveMux.Handle("GET /api/reset", http.HandlerFunc(apiCfg.reset))
	serveMux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(apiCfg.getChirp))
	serveMux.Handle("GET /api/chirps", http.HandlerFunc(apiCfg.getAllChirps))
	serveMux.Handle("POST /api/chirps", http.HandlerFunc(apiCfg.createChirp))
	serveMux.Handle("POST /api/login", http.HandlerFunc(apiCfg.login))
	serveMux.Handle("POST /api/users", http.HandlerFunc(apiCfg.createUser))
	serveMux.Handle("POST /api/refresh", http.HandlerFunc(apiCfg.refresh))
	serveMux.Handle("POST /api/revoke", http.HandlerFunc(apiCfg.revoke))
	serveMux.Handle("POST /api/polka/webhooks", http.HandlerFunc(apiCfg.webhooks))
	serveMux.Handle("PUT /api/users", http.HandlerFunc(apiCfg.updateUser))
	serveMux.Handle("DELETE /api/chirps/{chirpID}", http.HandlerFunc(apiCfg.deleteChirp))
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
