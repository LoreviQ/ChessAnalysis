package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db      *sql.DB
	queries map[string]string
	testDB  int
}

// NewConnection creates a new connection to the SQLite3 database
func NewConnection(test int) (*Database, error) {
	// Define the database path and schema path
	databasePath := "database.db"
	schemaPath := filepath.Join("sql", "schema.sql")
	if test != 0 {
		databasePath = fmt.Sprintf("test_database_%d.db", test)
	}

	// Open a connection to the SQLite3 database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	// Read the schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	// Execute the SQL commands in the schema.sql file
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	preparedQueries := map[string]string{
		"INSERT_MOVES":       "INSERT INTO moves (game_id, move_data) VALUES (?, ?)",
		"INSERT_GAME":        "INSERT INTO games (chessdotcom_id, playerIsWhite) VALUES (?, ?) RETURNING id",
		"GET_LATEST_GAME_ID": "SELECT id FROM games WHERE chessdotcom_id = ? ORDER BY created_at DESC LIMIT 1",
		"GET_LATEST_MOVES":   "SELECT id, move_data, scores, depth FROM moves WHERE game_id = ? ORDER BY created_at DESC LIMIT 1",
		"GET_GAMES":          "SELECT id, created_at, chessdotcom_id, playerIsWhite FROM games",
		"UPDATE_EVAL":        "UPDATE moves SET scores = ?, depth = ? WHERE id = ?",
	}

	return &Database{
		db:      db,
		queries: preparedQueries,
		testDB:  test,
	}, nil
}

// Close closes the connection to the SQLite3 database
func (d Database) Close() {
	d.db.Close()
	// Cleanup the test database
	if d.testDB != 0 {
		os.Remove(fmt.Sprintf("test_database_%d.db", d.testDB))
	}
}
