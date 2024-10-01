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
func (m Move) longAlgebraicNotation() (string, error) {
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
func (m Move) shortAlgebraicNotation() (string, error) {
	var piece, fromFile, capture, promotion, checkStatus string
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
			fromFile = string(m.FromFile)
		}
	}
	if m.Promotion != 0 {
		promotion = fmt.Sprintf("=%c", m.Promotion)
	}
	if m.CheckStatus != 0 {
		checkStatus = string(m.CheckStatus)
	}
	return fmt.Sprintf(
		"%s%s%s%c%d%s%s",
		piece,
		fromFile,
		capture,
		m.ToFile,
		m.ToRank,
		promotion,
		checkStatus,
	), nil
}
