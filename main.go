package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/janeturovskaya/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
	polkaKey       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		tokenSecret:    tokenSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	handlerFileServer := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", handlerFileServer)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUsers)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpById)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateEmailPassword)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaEvent)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
