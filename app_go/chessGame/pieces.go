package chessGame

import "fmt"

type piece struct {
	pType  pieceType
	color  string
	active bool
	rank   int
	file   string
}

type pieceType int8

const (
	NoPieceType pieceType = iota
	King
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

func (p *piece) getSymbol() (string, error) {
	switch p.pType {
	case King:
		return "K", nil
	case Queen:
		return "Q", nil
	case Rook:
		return "R", nil
	case Bishop:
		return "B", nil
	case Knight:
		return "N", nil
	case Pawn:
		return "", nil
	default:
		return "", fmt.Errorf("Invalid piece type")
	}
}

func (p *piece) getPrintable() (string, error) {
	switch p.pType {
	case Pawn:
		if p.color == "white" {
			return "♟", nil
		}
		return "♙", nil
	case King:
		if p.color == "white" {
			return "♚", nil
		}
		return "♔", nil
	case Queen:
		if p.color == "white" {
			return "♛", nil
		}
		return "♕", nil
	case Rook:
		if p.color == "white" {
			return "♜", nil
		}
		return "♖", nil
	case Bishop:
		if p.color == "white" {
			return "♝", nil
		}
		return "♗", nil
	case Knight:
		if p.color == "white" {
			return "♞", nil
		}
		return "♘", nil
	default:
		return "", fmt.Errorf("Invalid piece type")
	}
}
