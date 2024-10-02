package game

import (
	"errors"
	"fmt"
)

var (
	ErrNoPieceAtSquare = errors.New("no piece at square")
	ErrSquareOccupied  = errors.New("square is occupied")
	ErrPieceNotPawn    = errors.New("piece is not a pawn")
	ErrInvalidRank     = errors.New("invalid rank")
	ErrInvalidFile     = errors.New("invalid file")
	ErrPieceNotFound   = errors.New("piece not found")
)

type Board struct {
	Squares  [8][8]*Piece
	captured []*Piece
}

// NewBoard creates a new board with the initial game state
func NewBoard() *Board {
	b := Board{}
	b.setup_game()
	return &b
}

// Create a custom board with the given squares
func CustomBoard(squares [8][8]*Piece) *Board {
	return &Board{
		Squares: squares,
	}
}

// Set up the board with the initial game state
func (b *Board) setup_game() {
	b.Squares[0] = [8]*Piece{
		{PieceType: Rook, Color: "white", Active: true},
		{PieceType: Knight, Color: "white", Active: true},
		{PieceType: Bishop, Color: "white", Active: true},
		{PieceType: Queen, Color: "white", Active: true},
		{PieceType: King, Color: "white", Active: true},
		{PieceType: Bishop, Color: "white", Active: true},
		{PieceType: Knight, Color: "white", Active: true},
		{PieceType: Rook, Color: "white", Active: true},
	}
	b.Squares[7] = [8]*Piece{
		{PieceType: Rook, Color: "black", Active: true},
		{PieceType: Knight, Color: "black", Active: true},
		{PieceType: Bishop, Color: "black", Active: true},
		{PieceType: Queen, Color: "black", Active: true},
		{PieceType: King, Color: "black", Active: true},
		{PieceType: Bishop, Color: "black", Active: true},
		{PieceType: Knight, Color: "black", Active: true},
		{PieceType: Rook, Color: "black", Active: true},
	}
	for i := 0; i < 8; i++ {
		b.Squares[1][i] = &Piece{PieceType: Pawn, Color: "white", Active: true}
	}
	for i := 0; i < 8; i++ {
		b.Squares[6][i] = &Piece{PieceType: Pawn, Color: "black", Active: true}
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
				board += fmt.Sprintf("%c ", printable)
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
func (b *Board) GetPieceAtSquare(file rune, rank int) (*Piece, error) {
	if rank < 1 || rank > 8 {
		return nil, ErrInvalidRank
	}
	if fileToInt(file) < 1 || fileToInt(file) > 8 {
		return nil, ErrInvalidFile
	}
	return b.Squares[rank-1][fileToInt(file)-1], nil
}

// Move a piece from one square to another
// Doesn't check if the move is valid only if the square is occupied
func (b *Board) MovePiece(move Move) error {
	fromPiece, err := b.GetPieceAtSquare(move.FromFile, move.FromRank)
	if err != nil {
		return err
	}
	if fromPiece == nil {
		return ErrNoPieceAtSquare
	}
	toPiece, err := b.GetPieceAtSquare(move.ToFile, move.ToRank)
	if err != nil {
		return err
	}
	if move.Capture == 'x' {
		if toPiece == nil {
			return ErrNoPieceAtSquare
		}
		if toPiece.Color == fromPiece.Color {
			return ErrSquareOccupied
		}
		toPiece.Active = false
		b.captured = append(b.captured, toPiece)
	} else {
		if toPiece != nil {
			return ErrSquareOccupied
		}
	}
	b.Squares[move.ToRank-1][fileToInt(move.ToFile)-1] = fromPiece
	b.Squares[move.FromRank-1][fileToInt(move.FromFile-1)] = nil
	// TODO promotion
	return nil
}

// Promote a pawn to another piece type
func (b *Board) PromotePawn(file rune, rank int, pType PieceType) error {
	p, err := b.GetPieceAtSquare(file, rank)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrNoPieceAtSquare
	}
	if p.PieceType != Pawn {
		return ErrPieceNotPawn
	}
	p.PieceType = pType
	return nil
}

// Get the pieces captured by a given colour
// e.g. twken by white returns black pieces
func (b *Board) GetCapturedByColour(color string) []*Piece {
	var CapturedByColour []*Piece
	for _, p := range b.captured {
		if p.Color != color {
			CapturedByColour = append(CapturedByColour, p)
		}
	}
	return CapturedByColour
}

// converts 1-8 to a-h
func intToFile(i int) rune {
	return rune(i+'a') - 1
}

// converts a-h to 1-8
func fileToInt(r rune) int {
	return int(r-'a') + 1
}

func (b *Board) GetLocation(piece *Piece) (rune, int, error) {
	for i, row := range b.Squares {
		for j, p := range row {
			if p == piece {
				return intToFile(j + 1), i + 1, nil
			}
		}
	}
	return 0, 0, ErrPieceNotFound
}
