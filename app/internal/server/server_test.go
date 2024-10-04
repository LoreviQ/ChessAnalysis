package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LoreviQ/ChessAnalysis/app/internal/database"
)

func TestNewServer(t *testing.T) {
	// Create a new server
	srv, _ := NewServer(nil)

	// Check that the server is not nil
	if srv == nil {
		t.Error("NewServer() returned nil")
	}
}

func TestReadinessEndpoint(t *testing.T) {
	// Change the working directory to the root of the project
	restore := changeDirectoryToRoot()
	defer restore()

	// Create db connection
	db, err := database.NewConnection(4)
	if err != nil {
		t.Errorf("Error creating database connection: %v", err)
	}
	defer db.Close()

	// Create a new server
	srv, cfg := NewServer(db)
	go srv.ListenAndServe()
	defer srv.Close()
	url := cfg.url.String()

	// wait one second for the server to start
	time.Sleep(1 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/readiness", url))
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

// First tests the postMoves method by inserting a list of moves into the database
// and then tests the getLatestMoves method by retrieving the moves from the database
func TestPostGetMoves(t *testing.T) {
	// Change the working directory to the root of the project
	restore := changeDirectoryToRoot()
	defer restore()

	// Create db connection
	db, err := database.NewConnection(5)
	if err != nil {
		t.Errorf("Error creating database connection: %v", err)
	}
	defer db.Close()

	// Create a new server
	srv, cfg := NewServer(db)
	go srv.ListenAndServe()
	defer srv.Close()
	url := cfg.url.String()

	waitForServerToStart(url)

	movesToInsert := []string{
		"1", "e4", "g6", "2", "d4", "Bg7", "3", "e5", "d6",
		"4", "exd6", "Qxd6", "5", "Nc3", "Bxd4", "6", "Be3", "Bxc3+",
		"7", "bxc3", "h5", "8", "f4", "Qxd1+", "9", "Rxd1", "Nd7",
		"10", "Bc4", "e5", "11", "fxe5", "Nxe5", "12", "Bd4", "Nxc4",
		"13", "Bxh8", "Nb6", "14", "Nf3", "Nd7", "15", "O-O", "Nh6",
		"16", "Bf6", "Nxf6", "17", "Ng5", "Nd7", "18", "Nxf7", "Nxf7",
		"19", "Rd3", "Ke7",
	}
	expectedMoves := []string{
		"e2e4", "g7g6", "d2d4", "Bf8g7", "e4e5", "d7d6", "e5xd6", "Qd8xd6",
		"Nb1c3", "Bg7xd4", "Bc1e3", "Bd4xc3+", "b2xc3", "h7h5", "f2f4",
		"Qd6xd1+", "Ra1xd1", "Nb8d7", "Bf1c4", "e7e5", "f4xe5", "Nd7xe5",
		"Be3d4", "Ne5xc4", "Bd4xh8", "Nc4b6", "Ng1f3", "Nb6d7", "O-O", "Ng8h6",
		"Bh8f6", "Nd7xf6", "Nf3g5", "Nf6d7", "Ng5xf7", "Nh6xf7", "Rd1d3", "Ke8e7",
	}

	// Insert moves
	body, err := json.Marshal(map[string][]string{"moves": movesToInsert})
	if err != nil {
		t.Errorf("Error marshalling request body: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/games/123456/moves", url),
		"application/json",
		strings.NewReader(string(body)),
	)
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Get moves
	resp, err = http.Get(fmt.Sprintf("%s/games/123456/moves/latest", url))
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response getLatestMoveResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	resp.Body.Close()

	// Check that the moves are correct
	if len(response.Moves) != len(expectedMoves) {
		t.Errorf("Expected %d moves, got %d", len(expectedMoves), len(response.Moves))
	}
	for i := range response.Moves {
		if response.Moves[i] != expectedMoves[i] {
			t.Errorf("Expected move %s, got %s", expectedMoves[i], response.Moves[i])
		}
	}
}

// blocking function that waits for the server to start
func waitForServerToStart(url string) {
	for {
		resp, err := http.Get(fmt.Sprintf("%s/readiness", url))
		if err != nil {
			continue
		}
		if resp != nil {
			if resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				break
			}
			resp.Body.Close()
		}
	}
}

// changeDirectoryToRoot changes the working directory to the root of the project
// Returns a function that can be used to change back to the original directory
func changeDirectoryToRoot() func() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	// Change the working directory to the root of the project
	projectRoot := filepath.Join(cwd, "../..")
	err = os.Chdir(projectRoot)
	if err != nil {
		log.Fatalf("failed to change working directory: %v", err)
	}

	// Return a function that can be used to change back to the original directory
	return func() {
		err := os.Chdir(cwd)
		if err != nil {
			log.Fatalf("failed to change working directory: %v", err)
		}
	}
}
