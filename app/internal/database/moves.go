package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/LoreviQ/ChessAnalysis/app/internal/game"
)

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

// GetMovesByChessdotcomID returns the latest moves of a game with the given chess.com id
func (d Database) GetMovesByChessdotcomID(chessdotcomID string) ([]string, error) {
	var gameID int
	err := d.db.QueryRow(d.queries["GET_LATEST_GAME_ID"], chessdotcomID).Scan(&gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}
	return d.GetMovesByID(gameID)
}

// GetMoves returns the latest moves of a game with the given id
func (d Database) GetMovesByID(id int) ([]string, error) {
	var moves string
	err := d.db.QueryRow(d.queries["GET_LATEST_MOVES"], id).Scan(&moves)
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
