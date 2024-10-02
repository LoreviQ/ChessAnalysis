package game

import (
	"errors"
	"fmt"
)

var (
	ErrNotEnoughInfo = errors.New("not enough information to create long algebraic notation")
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
