package game

import (
	"regexp"
)

type Game struct {
	Board       *Board
	Turn        string
	MoveHistory []Move
	CanCastle   map[string]map[string]bool
}

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

// Create a new game
func NewGame() *Game {
	return &Game{
		Board:       NewBoard(),
		Turn:        "white",
		MoveHistory: []Move{},
		CanCastle: map[string]map[string]bool{
			"white": {"short": true, "long": true},
			"black": {"short": true, "long": true},
		},
	}
}

// Takes a move string in algebraic notation,
// checks if it is valid and moves the piece
func (g *Game) MovePiece(moveStr string) error {
	_, err := parseRegex(moveStr)
	if err != nil {
		return err
	}
	g.GetPossibleMoves()
	return nil
}

// Get all possible moves for the current player
func (g *Game) GetPossibleMoves() map[Move]*Piece {
	possibleMoves := map[Move]*Piece{}
	for _, row := range g.Board.Squares {
		for _, p := range row {
			if p != nil && p.Color == g.Turn {
				pieceMoves := p.GetPossibleMoves(g)
				for _, move := range pieceMoves {
					if possibleMoves[move] == nil {
						possibleMoves[move] = p
					} else {
						possibleMoves[move] = nil
						handleDuplicateMove(move)
					}
				}
			}
		}
	}
	return possibleMoves
}

// Handle a duplicate move
func handleDuplicateMove(move Move) {
	// Handle duplicate move
}

// Parse a move string in algebraic notation
func parseRegex(moveStr string) (Move, error) {
	// Regex to parse move string
	pattern := `^([NBRQK])?([a-h])?([1-8])?(x)?([a-h])([1-8])(=[NBRQK])?(\+|#)?$|^O-O(-O)?$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Move{}, err
	}
	matches := re.FindStringSubmatch(moveStr)
	if matches == nil {
		return Move{}, err
	}
	move := Move{}
	if matches[0] == "O-O" {
		move.Castle = "short"
	} else if matches[0] == "O-O-O" {
		move.Castle = "long"
	} else {
		if matches[1] != "" {
			move.Piece = rune(matches[1][0])
		}
		if matches[2] != "" {
			move.FromFile = rune(matches[2][0])
		}
		if matches[3] != "" {
			move.FromRank = int(matches[3][0] - '0')
		}
		if matches[4] != "" {
			move.Capture = rune(matches[4][0])
		}
		move.ToFile = rune(matches[5][0])
		move.ToRank = int(matches[6][0] - '0')
		if matches[7] != "" {
			move.Promotion = rune(matches[7][1])
		}
		if matches[8] != "" {
			move.CheckStatus = rune(matches[8][0])
		}
	}
	return move, nil
}
