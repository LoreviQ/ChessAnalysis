package game

type Board struct {
	squares [8][8]*piece
}

func NewBoard() Board {
	b := Board{}
	b.setup_game()
	return b
}

// Set up the board with the initial game state
func (b *Board) setup_game() {
	b.squares[0] = [8]*piece{
		{pType: Rook, color: "black", active: true},
		{pType: Knight, color: "black", active: true},
		{pType: Bishop, color: "black", active: true},
		{pType: Queen, color: "black", active: true},
		{pType: King, color: "black", active: true},
		{pType: Bishop, color: "black", active: true},
		{pType: Knight, color: "black", active: true},
		{pType: Rook, color: "black", active: true},
	}
	b.squares[7] = [8]*piece{
		{pType: Rook, color: "white", active: true},
		{pType: Knight, color: "white", active: true},
		{pType: Bishop, color: "white", active: true},
		{pType: Queen, color: "white", active: true},
		{pType: King, color: "white", active: true},
		{pType: Bishop, color: "white", active: true},
		{pType: Knight, color: "white", active: true},
		{pType: Rook, color: "white", active: true},
	}
	for i := 0; i < 8; i++ {
		b.squares[1][i] = &piece{pType: Pawn, color: "black", active: true}
	}
	for i := 0; i < 8; i++ {
		b.squares[6][i] = &piece{pType: Pawn, color: "white", active: true}
	}

}

func (b *Board) PrintBoard() string {
	var board string
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := b.squares[i][j]
			if p != nil {
				printable, _ := p.getPrintable()
				board += printable
			} else {
				board += " "
			}
		}
		board += "\n"
	}
	return board
}

func (b *Board) getPieceAtSquare(rank int, file string) *piece {
	files := map[string]int{
		"a": 0,
		"b": 1,
		"c": 2,
		"d": 3,
		"e": 4,
		"f": 5,
		"g": 6,
		"h": 7,
	}
	return b.squares[rank][files[file]]
}
