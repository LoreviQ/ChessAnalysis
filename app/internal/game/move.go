package game

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrNotEnoughInfo = errors.New("not enough information to create long algebraic notation")
	ErrAmbiguousMove = errors.New("ambiguous move, not enough information to determine the move")
)

type Move struct {
	Piece       rune
	FromFile    rune
	FromRank    int
	Capture     rune
	ToFile      rune
	ToRank      int
	Promotion   rune
	CheckStatus rune
	Castle      string
}

// Returns the long algebraic notation of the move
func (m Move) LongAlgebraicNotation() (string, error) {
	var piece, capture, promotion, checkStatus string
	if m.Castle == "short" {
		return "O-O", nil
	} else if m.Castle == "long" {
		return "O-O-O", nil
	}
	if m.FromFile == 0 || m.FromRank == 0 || m.ToFile == 0 || m.ToRank == 0 {
		return "", ErrNotEnoughInfo
	}
	if m.Piece != 0 {
		piece = string(m.Piece)
	}
	if m.Capture != 0 {
		capture = "x"
	}
	if m.Promotion != 0 {
		promotion = fmt.Sprintf("=%c", m.Promotion)
	}
	if m.CheckStatus != 0 {
		checkStatus = string(m.CheckStatus)
	}
	return fmt.Sprintf(
		"%s%c%d%s%c%d%s%s",
		piece,
		m.FromFile,
		m.FromRank,
		capture,
		m.ToFile,
		m.ToRank,
		promotion,
		checkStatus,
	), nil
}

// Returns the short algebraic notation of the move
// Does not consider when multiple pieces of the same type can move to the same square
func (m Move) ShortAlgebraicNotation(includeFile, includeRank bool) (string, error) {
	var piece, fromFile, fromRank, capture, promotion, checkStatus string
	if m.Castle == "short" {
		return "O-O", nil
	} else if m.Castle == "long" {
		return "O-O-O", nil
	}
	if m.ToFile == 0 || m.ToRank == 0 {
		return "", ErrNotEnoughInfo
	}
	if m.Piece != 0 {
		piece = string(m.Piece)
	}
	if m.Capture != 0 {
		capture = "x"
		if m.Piece == 0 { // pawn captures always display file
			includeFile = true
		}
	}
	if m.Promotion != 0 {
		promotion = fmt.Sprintf("=%c", m.Promotion)
	}
	if m.CheckStatus != 0 {
		checkStatus = string(m.CheckStatus)
	}
	if includeFile {
		if m.FromFile == 0 {
			return "", ErrNotEnoughInfo
		}
		fromFile = string(m.FromFile)
	}
	if includeRank {
		if m.FromRank == 0 {
			return "", ErrNotEnoughInfo
		}
		fromRank = fmt.Sprintf("%d", m.FromRank)
	}
	return fmt.Sprintf(
		"%s%s%s%s%c%d%s%s",
		piece,
		fromFile,
		fromRank,
		capture,
		m.ToFile,
		m.ToRank,
		promotion,
		checkStatus,
	), nil
}

// Returns the UCI notation of the move
func (m Move) UCInotation() (string, error) {
	var promotion string
	if m.FromRank == 0 {
		return "", ErrNotEnoughInfo
	}
	if m.Castle == "short" {
		return fmt.Sprintf("e%dg%d", m.FromRank, m.FromRank), nil
	} else if m.Castle == "long" {
		return fmt.Sprintf("e%dc%d", m.FromRank, m.FromRank), nil
	}
	if m.FromFile == 0 || m.ToFile == 0 || m.ToRank == 0 {
		return "", ErrNotEnoughInfo
	}
	if m.Promotion != 0 {
		promotion = string(unicode.ToLower(m.Promotion))
	}
	return fmt.Sprintf(
		"%c%d%c%d%s",
		m.FromFile,
		m.FromRank,
		m.ToFile,
		m.ToRank,
		promotion,
	), nil
}

// Get a list of strings represnting the short algebraic notation of provided moves
// Disambiguates between duplicate moves
func ConvertMovesToShortAlgebraicNotation(moves []Move) (map[string]Move, error) {
	notationToMove := map[string]Move{}
	for _, move := range moves {
		shortAlgebraicNotation, err := move.ShortAlgebraicNotation(false, false)
		if err != nil {
			return nil, err
		}
		otherMove, ok := notationToMove[shortAlgebraicNotation]
		if ok {
			// Disambiguate between duplicate moves
			var includeFile, includeRank bool
			if otherMove.FromFile == move.FromFile {
				// If the pieces are on the same file, disambiguate by rank
				includeFile, includeRank = false, true
			} else {
				// Otherwise, disambiguate by file
				includeFile, includeRank = true, false
			}
			newShortAlgebraicNotation, err := move.ShortAlgebraicNotation(includeFile, includeRank)
			if err != nil {
				return nil, err
			}
			otherShortAlgebraicNotation, err := otherMove.ShortAlgebraicNotation(includeFile, includeRank)
			if err != nil {
				return nil, err
			}
			notationToMove[newShortAlgebraicNotation] = move
			notationToMove[otherShortAlgebraicNotation] = otherMove
			delete(notationToMove, shortAlgebraicNotation)
		} else {
			notationToMove[shortAlgebraicNotation] = move
		}
	}
	return notationToMove, nil
}

func ConvertMovesToLongAlgebraicNotation(moves []Move) []string {
	notations := []string{}
	for _, move := range moves {
		longAlgebraicNotation, err := move.LongAlgebraicNotation()
		if err != nil {
			return []string{}
		}
		notations = append(notations, longAlgebraicNotation)
	}
	return notations
}

// Takes a move and a list of possible moves and returns the corresponding move from the list
// If the move is ambigious it returns an error
func getCorrespondingMove(move Move, possibleMoves []Move) (Move, error) {
	// Create a new slice to hold the filtered moves
	for _, filterType := range []string{"castle", "mandatory", "file", "rank", "promotion", "check"} {
		possibleMoves = filterMoves(move, possibleMoves, filterType)
		if len(possibleMoves) == 1 {
			return possibleMoves[0], nil
		}
		if len(possibleMoves) == 0 {
			return Move{}, ErrInvalidMove
		}
	}
	return Move{}, ErrAmbiguousMove
}

func filterMoves(move Move, moves []Move, filterType string) []Move {
	var filteredMoves []Move
	switch filterType {
	case "mandatory":
		for _, m := range moves {
			if m.Piece == move.Piece &&
				m.ToFile == move.ToFile &&
				m.ToRank == move.ToRank &&
				m.Capture == move.Capture {
				filteredMoves = append(filteredMoves, m)
			}
		}
	case "file":
		if move.FromFile == 0 {
			return moves
		}
		for _, m := range moves {
			if m.FromFile == move.FromFile {
				filteredMoves = append(filteredMoves, m)
			}
		}
	case "rank":
		if move.FromRank == 0 {
			return moves
		}
		for _, m := range moves {
			if m.FromRank == move.FromRank {
				filteredMoves = append(filteredMoves, m)
			}
		}
	case "promotion":
		if move.Promotion == 0 {
			return moves
		}
		for _, m := range moves {
			if m.Promotion == move.Promotion {
				filteredMoves = append(filteredMoves, m)
			}
		}
	case "check":
		if move.CheckStatus == 0 {
			return moves
		}
		for _, m := range moves {
			if m.CheckStatus == move.CheckStatus {
				filteredMoves = append(filteredMoves, m)
			}
		}
	case "castle":
		if move.Castle == "" {
			return moves
		}
		for _, m := range moves {
			if m.Castle == move.Castle {
				filteredMoves = append(filteredMoves, m)
			}
		}
	}
	return filteredMoves
}
