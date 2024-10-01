package game

import (
	"regexp"
	"strconv"
)

type Game struct {
	Board       Board
	Turn        string
	MoveHistory []string
	CanCastle   map[string]map[string]bool
}

type Move struct {
	Piece       string
	FromFile    string
	FromRank    int
	Capture     string
	ToFile      string
	ToRank      int
	Promotion   string
	CheckStatus string
	Castle      string
}

func NewGame() Game {
	return Game{
		Board:       NewBoard(),
		Turn:        "white",
		MoveHistory: []string{},
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
func (g *Game) GetPossibleMoves() []Move {
	possibleMoves := []Move{}
	for rank, row := range g.Board.Squares {
		for file, p := range row {
			if p != nil && p.color == g.Turn {
				possibleMoves = append(possibleMoves, p.GetPossibleMoves(g, intToFile(file), rank)...)
			}
		}
	}
	return possibleMoves
}

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
		move.Piece = matches[1]
		move.FromFile = matches[2]
		move.FromRank, _ = strconv.Atoi(matches[3])
		move.Capture = matches[4]
		move.ToFile = matches[5]
		move.ToRank, _ = strconv.Atoi(matches[6])
		move.Promotion = matches[7]
		move.CheckStatus = matches[8]
	}

	return move, nil
}
