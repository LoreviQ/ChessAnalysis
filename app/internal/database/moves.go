package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/LoreviQ/ChessAnalysis/app/internal/eval"
	"github.com/LoreviQ/ChessAnalysis/app/internal/game"
)

var (
	ErrNoMoves = errors.New("game has no moves")
)

type Move struct {
	ID     int
	Moves  []string
	Scores []string
	Depth  int
}

// InsertMoves inserts a list of moves into the database
func (d Database) InsertMoves(moves []string, chessdotcomID string, playerIsWhite bool) error {
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
		err = d.db.QueryRow(d.queries["INSERT_GAME"], chessdotcomID_NullString, playerIsWhite).Scan(&gameID)
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
func (d Database) GetMovesByChessdotcomID(chessdotcomID string) (*Move, error) {
	var gameID int
	err := d.db.QueryRow(d.queries["GET_LATEST_GAME_ID"], chessdotcomID).Scan(&gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}
	return d.GetMovesByID(gameID)
}

// GetMoves returns the latest moves of a game with the given id
func (d Database) GetMovesByID(id int) (*Move, error) {
	var moves string
	var scores sql.NullString
	var moves_id int
	var depth sql.NullInt64
	err := d.db.QueryRow(d.queries["GET_LATEST_MOVES"], id).Scan(&moves_id, &moves, &scores, &depth)
	if err != nil {
		return nil, ErrNoMoves
	}
	scoresOut := []string{}
	if scores.Valid {
		scoresOut = strings.Split(scores.String, " ")
	}
	depthOut := 0
	if depth.Valid {
		depthOut = int(depth.Int64)
	}
	return &Move{
		ID:     moves_id,
		Moves:  strings.Split(moves, " "),
		Scores: scoresOut,
		Depth:  depthOut,
	}, nil
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

// UpdateEval updates the evaluation of a move in the database
func (d Database) UpdateEval(moveID int, evalss [][]*eval.MoveEval) error {
	scores := []string{}
	depth := 0
	for _, evals := range evalss {
		e := eval.GetEvalNum(evals, 1)
		if e == nil {
			continue
		}
		if e.Depth > depth {
			depth = e.Depth
		}
		if e.Mate {
			scores = append(scores, fmt.Sprintf("M%d", e.MateIn))
		} else {
			scores = append(scores, fmt.Sprintf("%d", e.Score))
		}
	}
	scoresStr := strings.Join(scores, " ")
	_, err := d.db.Exec(d.queries["UPDATE_EVAL"], scoresStr, depth, moveID)
	return err
}
