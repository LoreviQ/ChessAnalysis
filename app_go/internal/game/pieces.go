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
// K, Q, R, B, N, zeroval for Pawn
func (p *piece) getSymbol() (rune, error) {
	switch p.pType {
	case King:
		return 'K', nil
	case Queen:
		return 'Q', nil
	case Rook:
		return 'R', nil
	case Bishop:
		return 'B', nil
	case Knight:
		return 'N', nil
	case Pawn:
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid piece type")
	}
}

// Returns the symbol of the piece
// ♔, ♕, ♖, ♗, ♘, ♙, ♚, ♛, ♜, ♝, ♞, ♟
func (p *piece) getPrintable() (rune, error) {
	switch p.pType {
	case Pawn:
		if p.color == "white" {
			return '♟', nil
		}
		return '♙', nil
	case King:
		if p.color == "white" {
			return '♚', nil
		}
		return '♔', nil
	case Queen:
		if p.color == "white" {
			return '♛', nil
		}
		return '♕', nil
	case Rook:
		if p.color == "white" {
			return '♜', nil
		}
		return '♖', nil
	case Bishop:
		if p.color == "white" {
			return '♝', nil
		}
		return '♗', nil
	case Knight:
		if p.color == "white" {
			return '♞', nil
		}
		return '♘', nil
	default:
		return 0, fmt.Errorf("invalid piece type")
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
func (p *piece) GetPossibleMoves(g *Game, fromFile rune, fromRank int) []Move {
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
func (p *piece) getPawnMoves(g *Game, fromFile rune, fromRank int) []Move {
	moves := []Move{}
	direction := p.getDirection()
	// Forward one square
	toRank := fromRank + direction
	forward_piece, err := g.Board.GetPieceAtSquare(fromFile, toRank)
	if err == nil && forward_piece == nil {
		moves = append(moves, Move{
			FromFile: fromFile,
			FromRank: fromRank,
			ToFile:   fromFile,
			ToRank:   toRank,
		})
		// Forward two squares
		if (direction == 1 && fromRank == 2) || (direction == -1 && fromRank == 7) {
			toRank = fromRank + 2*direction
			forward_piece, err := g.Board.GetPieceAtSquare(fromFile, toRank)
			if err == nil && forward_piece == nil {
				moves = append(moves, Move{
					FromFile: fromFile,
					FromRank: fromRank,
					ToFile:   fromFile,
					ToRank:   toRank,
				})
			}
		}
	}

	return moves
}

func (p *piece) getKingMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getQueenMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getRookMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getBishopMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

func (p *piece) getKnightMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}
