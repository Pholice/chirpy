package main

import (
	"net/http"
)

type apiConfig struct {
	fileserverHits int
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
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	serveMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	serveMux.Handle("GET /api/healthz", http.HandlerFunc(apiCfg.healthz))
	serveMux.Handle("GET /admin/metrics", http.HandlerFunc(apiCfg.metrics))
	serveMux.Handle("GET /api/reset", http.HandlerFunc(apiCfg.reset))
	serveMux.Handle("POST /api/chirps", http.HandlerFunc(apiCfg.createChirp))

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
