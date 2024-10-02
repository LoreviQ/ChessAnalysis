package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LoreviQ/ChessAnalysis/app_go/internal/game"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db      *sql.DB
	queries map[string]string
}

func NewConnection() Database {
	// Get file paths
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get the current file path")
	}
	dir := filepath.Dir(filename)
	databasePath := filepath.Join(dir, "database.db")
	schemaPath := filepath.Join(dir, "sql", "schema.sql")

	// Open a connection to the SQLite3 database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		panic(err)
	}

	// Read the schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL commands in the schema.sql file
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	preparedQueries := map[string]string{
		"INSERT_MOVES":       "INSERT INTO moves (game_id, move_data) VALUES (?, ?)",
		"INSERT_GAME":        "INSERT INTO games (chessdotcom_id) VALUES (?) RETURNING id",
		"GET_LATEST_GAME_ID": "SELECT id FROM games WHERE chessdotcom_id = ? ORDER BY created_at DESC LIMIT 1",
	}

	return Database{db: db, queries: preparedQueries}
}

func (d Database) Close() {
	d.db.Close()
}

func (d Database) InsertMoves(moves []string, chessdotcomID string) error {
	chessdotcomID_NullString := sql.NullString{String: chessdotcomID, Valid: chessdotcomID != ""}
	standardizedMoves, err := standardizeMoves(moves)
	if err != nil {
		return err
	}

	// Get game id of the latest game with the given chess.com id
	var gameID int
	err = d.db.QueryRow(d.queries["GET_LATEST_GAME_ID"], chessdotcomID_NullString).Scan(&gameID)
	if err != nil {
		return err
	}

	// If no game with the given chess.com id exists, create a new game
	if gameID == 0 {
		err = d.db.QueryRow(d.queries["INSERT_GAME"], chessdotcomID_NullString).Scan(&gameID)
		if err != nil {
			return err
		}
	}

	// Insert the moves into the database
	_, err = d.db.Exec(d.queries["INSERT_MOVES"], gameID, standardizedMoves)
	return err
}

// Converts moves to the format used in the database
func standardizeMoves(moves []string) ([]string, error) {
	// remove turn numbers
	standardizedMoves := []string{}
	for i, move := range moves {
		if i%3 != 0 {
			standardizedMoves = append(standardizedMoves, move)
		}
	}
	standardizedMoves, err := game.ConvertNotation(standardizedMoves)
	if err != nil {
		return nil, err
	}
	return standardizedMoves, nil
}
