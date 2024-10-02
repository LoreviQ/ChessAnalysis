package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	// Create a new server
	srv, _ := NewServer()

	// Check that the server is not nil
	if srv == nil {
		t.Error("NewServer() returned nil")
	}
}

func TestReadinessEndpoint(t *testing.T) {
	// Create a new server
	srv, cfg := NewServer()
	go srv.ListenAndServe()
	defer srv.Close()

	// wait one second for the server to start
	time.Sleep(1 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/readiness", cfg.url.String()))
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
