package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/database"
	"github.com/joho/godotenv"
)

type serverCfg struct {
	url *url.URL
	db  *database.Database
}

func NewServer(db *database.Database) (*http.Server, serverCfg) {
	cfg := setupCfg(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /readiness", cfg.getReadiness)
	mux.HandleFunc("POST /games/{id}/moves", cfg.postMoves)
	mux.HandleFunc("GET /games/{id}/moves/latest", cfg.getLatestMoves)
	return &http.Server{
		Addr:    cfg.url.Host,
		Handler: CorsMiddleware(mux),
	}, *cfg

}

func setupCfg(db *database.Database) *serverCfg {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Assuming default configuration - .env unreadable: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default port
	}
	return &serverCfg{
		url: &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%s", port),
		},
		db: db,
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
