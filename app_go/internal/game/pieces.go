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
func (p *piece) GetPossibleMoves(g *Game, fromFile string, fromRank int) []Move {
	if !p.active {
		return []Move{}
	}
	switch p.pType {
	case Pawn:
		return p.getPawnMoves(g, fromFile, fromRank)
	case King:
		return p.getKingMoves(g, fromFile, fromRank)
	case Queen:
		return p.getQueenMoves(g, fromFile, fromRank)
	case Rook:
		return p.getRookMoves(g, fromFile, fromRank)
	case Bishop:
		return p.getBishopMoves(g, fromFile, fromRank)
	case Knight:
		return p.getKnightMoves(g, fromFile, fromRank)
	default:
		return []Move{}
	}
}

// Returns the possible moves for a pawn
func (p *piece) getPawnMoves(g *Game, fromFile string, fromRank int) []Move {
	// Forward one square
	toRank := fromRank + p.getDirection()
	print(fromFile, toRank)
	return []Move{}
}

func (p *piece) getKingMoves(g *Game, fromFile string, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getQueenMoves(g *Game, fromFile string, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getRookMoves(g *Game, fromFile string, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getBishopMoves(g *Game, fromFile string, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getKnightMoves(g *Game, fromFile string, fromRank int) []Move {
	return []Move{}
}
