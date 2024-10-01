package game

import "fmt"

type piece struct {
	pType  pieceType
	color  string
	active bool
}

type pieceType int8

const (
	King = iota
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

// Returns the letter representing the piece
// K, Q, R, B, N, ""
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
		return "", fmt.Errorf("invalid piece type")
	}
}

// Returns the symbol of the piece
// ♔, ♕, ♖, ♗, ♘, ♙, ♚, ♛, ♜, ♝, ♞, ♟
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
		return "", fmt.Errorf("invalid piece type")
	}
}

// Returns the direction the piece moves
func (p *piece) getDirection() int {
	if p.color == "white" {
		return 1
	}
	return -1
}

// Returns the possible moves for the piece
func (p *piece) GetPossibleMoves(g *Game, file string, rank int) []Move {
	switch p.pType {
	case Pawn:
		return p.getPawnMoves(g, file, rank)
	case King:
		return p.getKingMoves(g, file, rank)
	case Queen:
		return p.getQueenMoves(g, file, rank)
	case Rook:
		return p.getRookMoves(g, file, rank)
	case Bishop:
		return p.getBishopMoves(g, file, rank)
	case Knight:
		return p.getKnightMoves(g, file, rank)
	default:
		return []Move{}
	}
}

// Returns the possible moves for a pawn
func (p *piece) getPawnMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}

func (p *piece) getKingMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}

func (p *piece) getQueenMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}

func (p *piece) getRookMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}

func (p *piece) getBishopMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}

func (p *piece) getKnightMoves(g *Game, file string, rank int) []Move {
	return []Move{}
}
