package game

import (
	"errors"
	"fmt"
)

var (
	ErrNoPieceAtSquare = errors.New("no piece at square")
	ErrSquareOccupied  = errors.New("square is occupied")
	ErrPieceNotPawn    = errors.New("piece is not a pawn")
)

type Board struct {
	Squares  [8][8]*piece
	captured []*piece
}

// NewBoard creates a new board with the initial game state
func NewBoard() Board {
	b := Board{}

	b.setup_game()
	return b
}

// Set up the board with the initial game state
func (b *Board) setup_game() {
	b.Squares[0] = [8]*piece{
		{pType: Rook, color: "white", active: true},
		{pType: Knight, color: "white", active: true},
		{pType: Bishop, color: "white", active: true},
		{pType: Queen, color: "white", active: true},
		{pType: King, color: "white", active: true},
		{pType: Bishop, color: "white", active: true},
		{pType: Knight, color: "white", active: true},
		{pType: Rook, color: "white", active: true},
	}
	b.Squares[7] = [8]*piece{
		{pType: Rook, color: "black", active: true},
		{pType: Knight, color: "black", active: true},
		{pType: Bishop, color: "black", active: true},
		{pType: Queen, color: "black", active: true},
		{pType: King, color: "black", active: true},
		{pType: Bishop, color: "black", active: true},
		{pType: Knight, color: "black", active: true},
		{pType: Rook, color: "black", active: true},
	}
	for i := 0; i < 8; i++ {
		b.Squares[1][i] = &piece{pType: Pawn, color: "white", active: true}
	}
	for i := 0; i < 8; i++ {
		b.Squares[6][i] = &piece{pType: Pawn, color: "black", active: true}
	}

}

// PrintBoard returns a string representation of the board
func (b *Board) PrintBoard() string {
	var board string
	for i := 7; i >= 0; i-- { // Iterate from 7 to 0
		board += fmt.Sprintf("%d ", i+1)
		for j := 0; j < 8; j++ {
			p := b.Squares[i][j]
			if p != nil {
				printable, _ := p.getPrintable()
				board += fmt.Sprintf("%s ", printable)
			} else {
				board += "  "
			}
		}
		board += "\n"
	}
	board += "  a b c d e f g h\n"
	return board
}

// Get the piece at a given square
func (b *Board) GetPieceAtSquare(file string, rank int) *piece {
	return b.Squares[rank-1][fileToInt(file)]
}

// Move a piece from one square to another
// Doesn't check if the move is valid only if the square is occupied
func (b *Board) MovePiece(fromFile string, fromRank int, toFile string, toRank int) error {
	fromPiece := b.GetPieceAtSquare(fromFile, fromRank)
	if fromPiece == nil {
		return ErrNoPieceAtSquare
	}
	toPiece := b.GetPieceAtSquare(toFile, toRank)
	if toPiece != nil {
		return ErrSquareOccupied
	}
	b.Squares[toRank-1][fileToInt(toFile)] = fromPiece
	b.Squares[fromRank-1][fileToInt(fromFile)] = nil
	return nil
}

// Capture a piece from a square
func (b *Board) CapturePiece(file string, rank int) error {
	p := b.GetPieceAtSquare(file, rank)
	if p == nil {
		return ErrNoPieceAtSquare
	}
	p.active = false
	b.captured = append(b.captured, p)
	b.Squares[rank-1][fileToInt(file)] = nil
	return nil
}

// Promote a pawn to another piece type
func (b *Board) PromotePawn(file string, rank int, pType pieceType) error {
	p := b.GetPieceAtSquare(file, rank)
	if p == nil {
		return ErrNoPieceAtSquare
	}
	if p.pType != Pawn {
		return ErrPieceNotPawn
	}
	p.pType = pType
	return nil
}

// Get the pieces captured by a given colour
// e.g. twken by white returns black pieces
func (b *Board) GetCapturedByColour(color string) []*piece {
	var CapturedByColour []*piece
	for _, p := range b.captured {
		if p.color != color {
			CapturedByColour = append(CapturedByColour, p)
		}
	}
	return CapturedByColour
}

// converts 1-8 to a-h
func intToFile(i int) string {
	return fmt.Sprintf("%c", i+96)
}

// converts a-h to 0-7
func fileToInt(r string) int {
	return int(r[0] - 97)
}
