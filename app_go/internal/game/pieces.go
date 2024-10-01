package game

import "fmt"

type Piece struct {
	PieceType PieceType
	Color     string
	Active    bool
}

type PieceType int8

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
func (p *Piece) getSymbol() (rune, error) {
	switch p.PieceType {
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
func (p *Piece) getPrintable() (rune, error) {
	switch p.PieceType {
	case Pawn:
		if p.Color == "white" {
			return '♟', nil
		}
		return '♙', nil
	case King:
		if p.Color == "white" {
			return '♚', nil
		}
		return '♔', nil
	case Queen:
		if p.Color == "white" {
			return '♛', nil
		}
		return '♕', nil
	case Rook:
		if p.Color == "white" {
			return '♜', nil
		}
		return '♖', nil
	case Bishop:
		if p.Color == "white" {
			return '♝', nil
		}
		return '♗', nil
	case Knight:
		if p.Color == "white" {
			return '♞', nil
		}
		return '♘', nil
	default:
		return 0, fmt.Errorf("invalid piece type")
	}
}

// Returns the direction the piece moves
func (p *Piece) getDirection() int {
	if p.Color == "white" {
		return 1
	}
	return -1
}

// Returns the possible moves for the piece
func (p *Piece) GetPossibleMoves(g *Game, fromFile rune, fromRank int) []Move {
	if !p.Active {
		return []Move{}
	}
	switch p.PieceType {
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
func (p *Piece) getPawnMoves(g *Game, fromFile rune, fromRank int) []Move {
	moves := []Move{}
	moves = append(moves, p.pawnForward(g, fromFile, fromRank)...)
	moves = append(moves, p.pawnDiagonally(g, fromFile, fromRank)...)
	return moves
}

func (p *Piece) pawnForward(g *Game, fromFile rune, fromRank int) []Move {
	// Forward one square
	moves := []Move{}
	direction := p.getDirection()
	toRank := fromRank + direction
	forwardPiece, err := g.Board.GetPieceAtSquare(fromFile, toRank)
	if err != nil || forwardPiece != nil {
		return moves
	}
	move := Move{
		FromFile: fromFile,
		FromRank: fromRank,
		ToFile:   fromFile,
		ToRank:   toRank,
	}
	moves = append(moves, listPawnPromotions(move)...)

	// Forward two squares
	if !((direction == 1 && fromRank == 2) || (direction == -1 && fromRank == 7)) {
		return moves
	}
	toRank = fromRank + 2*direction
	doubleForwardPiece, err := g.Board.GetPieceAtSquare(fromFile, toRank)
	if err != nil || doubleForwardPiece != nil {
		return moves
	}
	moves = append(moves, Move{
		FromFile: fromFile,
		FromRank: fromRank,
		ToFile:   fromFile,
		ToRank:   toRank,
	})
	return moves
}

func (p *Piece) pawnDiagonally(g *Game, fromFile rune, fromRank int) []Move {
	moves := []Move{}
	direction := p.getDirection()
	// Capture diagonally
	toRank := fromRank + direction
	for _, toFile := range []rune{fromFile - 1, fromFile + 1} {
		diagonalPiece, err := g.Board.GetPieceAtSquare(toFile, toRank)
		if err != nil {
			continue
		}
		if diagonalPiece != nil &&
			diagonalPiece.Color != p.Color {
			move := Move{
				FromFile: fromFile,
				FromRank: fromRank,
				Capture:  'x',
				ToFile:   toFile,
				ToRank:   toRank,
			}
			moves = append(moves, listPawnPromotions(move)...)
		}
		// En passant
		if !((fromRank == 5 && direction == 1) ||
			(fromRank == 4 && direction == -1)) {
			continue
		}
		requiredPreviousMove := Move{
			FromFile: toFile,
			FromRank: toRank + direction,
			ToFile:   toFile,
			ToRank:   toRank - direction,
		}
		previousMove := g.MoveHistory[len(g.MoveHistory)-1]
		if previousMove != requiredPreviousMove {
			continue
		}
		moves = append(moves, Move{
			FromFile: fromFile,
			FromRank: fromRank,
			Capture:  'x',
			ToFile:   toFile,
			ToRank:   toRank,
		})
	}
	return moves
}

func listPawnPromotions(move Move) []Move {
	moves := []Move{}
	if move.ToRank == 8 || move.ToRank == 1 {
		for _, promotion := range []rune{'N', 'B', 'R', 'Q'} {
			move.Promotion = promotion
			moves = append(moves, move)
		}
	} else {
		moves = append(moves, move)
	}
	return moves
}

func (p *Piece) getKingMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

func (p *Piece) getQueenMoves(g *Game, fromFile rune, fromRank int) []Move {
	moves, _ := p.getMovesInDirection(g, fromFile, fromRank, "both")
	return moves
}

func (p *Piece) getRookMoves(g *Game, fromFile rune, fromRank int) []Move {
	moves, _ := p.getMovesInDirection(g, fromFile, fromRank, "orthogonal")
	return moves
}

func (p *Piece) getBishopMoves(g *Game, fromFile rune, fromRank int) []Move {
	moves, _ := p.getMovesInDirection(g, fromFile, fromRank, "diagonal")
	return moves
}

func (p *Piece) getKnightMoves(g *Game, fromFile rune, fromRank int) []Move {
	return []Move{}
}

// Returns the possible moves for a piece moving orthogonally
func (p *Piece) getMovesInDirection(g *Game, fromFile rune, fromRank int, moveType string) ([]Move, error) {
	directions := map[string][]int{}
	if moveType != "orthogonal" && moveType != "diagonal" && moveType != "both" {
		return nil, fmt.Errorf("invalid move type")
	}
	if moveType == "orthogonal" || moveType == "both" {
		directions["forward"] = []int{0, 1}
		directions["backward"] = []int{0, -1}
		directions["left"] = []int{-1, 0}
		directions["right"] = []int{1, 0}
	}
	if moveType == "diagonal" || moveType == "both" {
		directions["forward-right"] = []int{1, 1}
		directions["forward-left"] = []int{-1, 1}
		directions["backward-right"] = []int{1, -1}
		directions["backward-left"] = []int{-1, -1}
	}
	moves := []Move{}
	for direction := range directions {
		fileStep, rankStep := rune(directions[direction][0]), directions[direction][1]
		toFile, toRank := fromFile, fromRank
		for {
			toFile += fileStep
			toRank += rankStep
			toPiece, err := g.Board.GetPieceAtSquare(toFile, toRank)
			if err != nil {
				break
			}
			if toPiece != nil {
				if toPiece.Color != p.Color {
					moves = append(moves, Move{
						FromFile: fromFile,
						FromRank: fromRank,
						Capture:  'x',
						ToFile:   toFile,
						ToRank:   toRank,
					})
				}
				break
			}
			moves = append(moves, Move{
				FromFile: fromFile,
				FromRank: fromRank,
				ToFile:   toFile,
				ToRank:   toRank,
			})
		}
	}
	return moves, nil
}
