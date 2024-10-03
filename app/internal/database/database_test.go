package database

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConnection(t *testing.T) {
	// Change the working directory to the root of the project
	restore := changeDirectoryToRoot()
	defer restore()

	db, err := NewConnection(true)
	if err != nil {
		t.Error(err)
	}
	if db.db == nil {
		t.Error("Database connection is nil")
	}
	db.Close()
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
