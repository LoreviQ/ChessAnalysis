package game

import "fmt"

type Board struct {
	Squares  [8][8]*piece
	captured []*piece
	files    map[string]int
}

// NewBoard creates a new board with the initial game state
func NewBoard() Board {
	b := Board{
		files: map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
			"d": 3,
			"e": 4,
			"f": 5,
			"g": 6,
			"h": 7,
		},
	}
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
				board += printable
			} else {
				board += "\u3000"
			}
		}
		board += "\n"
	}
	board += "  a b c d e f g h\n"
	return board
}

// Get the piece at a given square
func (b *Board) GetPieceAtSquare(file string, rank int) *piece {
	return b.Squares[rank-1][b.files[file]]
}

// Move a piece from one square to another
// Doesn't check if the move is valid only if the square is occupied
func (b *Board) MovePiece(fromFile string, fromRank int, toFile string, toRank int) error {
	fromPiece := b.GetPieceAtSquare(fromFile, fromRank)
	if fromPiece == nil {
		return fmt.Errorf("no piece at square")
	}
	toPiece := b.GetPieceAtSquare(toFile, toRank)
	if toPiece != nil {
		return fmt.Errorf("square is occupied")
	}
	b.Squares[toRank-1][b.files[toFile]] = fromPiece
	b.Squares[fromRank-1][b.files[fromFile]] = nil
	return nil
}

// Capture a piece from a square
func (b *Board) CapturePiece(file string, rank int) error {
	p := b.GetPieceAtSquare(file, rank)
	if p == nil {
		return fmt.Errorf("no piece at square")
	}
	p.active = false
	b.captured = append(b.captured, p)
	b.Squares[rank-1][b.files[file]] = nil
	return nil
}

// Promote a pawn to another piece type
func (b *Board) PromotePawn(file string, rank int, pType pieceType) error {
	p := b.GetPieceAtSquare(file, rank)
	if p == nil || p.pType != Pawn {
		return fmt.Errorf("no pawn at square")
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
