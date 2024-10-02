package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db      *sql.DB
	queries map[string]string
	testDB  bool
}

// NewConnection creates a new connection to the SQLite3 database
func NewConnection(test bool) (Database, error) {
	// Get file paths
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get the current file path")
	}
	dir := filepath.Dir(filename)
	databasePath := filepath.Join(dir, "database.db")
	schemaPath := filepath.Join(dir, "sql", "schema.sql")
	if test {
		databasePath = filepath.Join(dir, "test_database.db")
	}

	// Open a connection to the SQLite3 database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return Database{}, err
	}

	// Read the schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return Database{}, err
	}

	// Execute the SQL commands in the schema.sql file
	_, err = db.Exec(string(schema))
	if err != nil {
		return Database{}, err
	}

	preparedQueries := map[string]string{
		"INSERT_MOVES":       "INSERT INTO moves (game_id, move_data) VALUES (?, ?)",
		"INSERT_GAME":        "INSERT INTO games (chessdotcom_id) VALUES (?) RETURNING id",
		"GET_LATEST_GAME_ID": "SELECT id FROM games WHERE chessdotcom_id = ? ORDER BY created_at DESC LIMIT 1",
		"GET_LATEST_MOVES":   "SELECT move_data FROM moves WHERE game_id = ? ORDER BY created_at DESC LIMIT 1",
	}

	return Database{
		db:      db,
		queries: preparedQueries,
		testDB:  test,
	}, nil
}

// Close closes the connection to the SQLite3 database
func (d Database) Close() {
	d.db.Close()
	// Cleanup the test database
	if d.testDB {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("Unable to get the current file path")
		}
		dir := filepath.Dir(filename)
		databasePath := filepath.Join(dir, "test_database.db")
		os.Remove(databasePath)
	}
}

// InsertMoves inserts a list of moves into the database
func (d Database) InsertMoves(moves []string, chessdotcomID string) error {
	chessdotcomID_NullString := sql.NullString{String: chessdotcomID, Valid: chessdotcomID != ""}
	standardizedMoves, err := standardizeMoves(moves)
	if err != nil {
		return err
	}

	// Get game id of the latest game with the given chess.com id
	var gameID int
	err = d.db.QueryRow(d.queries["GET_LATEST_GAME_ID"], chessdotcomID_NullString).Scan(&gameID)
	if err == sql.ErrNoRows {
		// If no game with the given chess.com id exists, create a new game
		err = d.db.QueryRow(d.queries["INSERT_GAME"], chessdotcomID_NullString).Scan(&gameID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Insert the moves into the database
	_, err = d.db.Exec(d.queries["INSERT_MOVES"], gameID, standardizedMoves)
	return err
}

// GetMoves returns the latest moves of a game with the given chess.com id
func (d Database) GetMoves(chessdotcomID string) ([]string, error) {
	var gameID int
	err := d.db.QueryRow(d.queries["GET_LATEST_GAME_ID"], chessdotcomID).Scan(&gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}
	var moves string
	err = d.db.QueryRow(d.queries["GET_LATEST_MOVES"], gameID).Scan(&moves)
	if err != nil {
		return nil, errors.New("game has no moves")
	}
	return strings.Split(moves, " "), nil
}

// Converts moves to the format used in the database
//
// Expected input: ["1", "e4", "e5", "2", "Nf3", "Nc6", ...]
// Expected output: "e4 e5 Nf3 Nc6 ..."
func standardizeMoves(moves []string) (string, error) {
	// remove turn numbers
	standardizedMoves := []string{}
	for i, move := range moves {
		if i%3 != 0 {
			standardizedMoves = append(standardizedMoves, move)
		}
	}
	standardizedMoves, err := game.ConvertNotation(standardizedMoves)
	if err != nil {
		return "", err
	}
	moveString := strings.Join(standardizedMoves, " ")
	return moveString, nil
}
