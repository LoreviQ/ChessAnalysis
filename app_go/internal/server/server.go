package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type serverCfg struct {
	port string
}

func NewServer() *http.Server {
	cfg := setupCfg()
	mux := http.NewServeMux()
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.port),
		Handler: mux,
	}
}

func setupCfg() *serverCfg {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Assuming default configuration - .env unreadable: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default port
	}
	return &serverCfg{
		port: port,
	}
}
