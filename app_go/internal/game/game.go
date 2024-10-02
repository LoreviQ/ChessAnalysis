package game

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var ErrInvalidMove = errors.New("invalid move")

type Game struct {
	Board       *Board
	Turn        string
	MoveHistory []Move
}

// Create a new game
func NewGame() *Game {
	return &Game{
		Board:       NewBoard(),
		Turn:        "white",
		MoveHistory: []Move{},
	}
}

// Starts playing the chess game in the console
func (g *Game) Play() {
	for {
		fmt.Printf("\n%s", g.Board.PrintBoard())
		userInput := getUserInput(fmt.Sprintf("%s to move: ", g.Turn))
		args := strings.Split(strings.ToLower(userInput), " ")
		switch args[0] {
		case "help":
			fmt.Println("Type move in short algebraic notation to play it (e.g. e4)")
			fmt.Println("Type 'quit' to exit the game")
			fmt.Println("Type 'move_history' to see the move history")
			fmt.Println("      Add '--short' to see the move history in short algebraic notation")
			fmt.Println("Type 'possible_moves' to see all possible moves")
			continue
		case "quit":
			return
		case "move_history":
			printString := ""
			if len(args) > 1 && args[1] == "--short" {
				notations, err := ConvertMovesToShortAlgebraicNotation(g.MoveHistory)
				if err != nil {
					fmt.Println(err)
				}
				for notation := range notations {
					printString += notation + ", "
				}
			} else {
				notations := ConvertMovesToLongAlgebraicNotation(g.MoveHistory)
				for _, notation := range notations {
					printString += notation + ", "
				}
			}
			fmt.Println("Previous moves:")
			fmt.Println(printString[:len(printString)-2])
			continue
		case "possible_moves":
			g.logPossibleMoves()
			continue
		default:
			err := g.Move(userInput)
			if err == ErrInvalidMove {
				fmt.Println("Invalid move")
				g.logPossibleMoves()

			} else if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Log all possible moves for the current player
func (g *Game) logPossibleMoves() {
	possibleMoves := g.GetPossibleMoves()
	notations, err := ConvertMovesToShortAlgebraicNotation(possibleMoves)
	if err != nil {
		fmt.Println(err)
	}
	printString := ""
	for notation := range notations {
		printString += notation + ", "
	}
	fmt.Println("Possible moves:")
	fmt.Println(printString[:len(printString)-2])
}

// Takes a move string in algebraic notation,
// checks if it is valid and moves the piece
func (g *Game) Move(moveStr string) error {
	// Check if the move string is valid
	_, err := parseRegex(moveStr)
	if err != nil {
		return err
	}
	// Get all possible moves for the current player
	possibleMoves := g.GetPossibleMoves()
	notationToMove, err := ConvertMovesToShortAlgebraicNotation(possibleMoves)
	if err != nil {
		return err
	}
	// Find the move that produces the given notation and play it
	for notation := range notationToMove {
		if notation == moveStr {
			move := notationToMove[notation]
			err = g.Board.MovePiece(move)
			if err != nil {
				return err
			}
			g.MoveHistory = append(g.MoveHistory, move)
			g.changeTurn()
			return nil
		}
	}
	return ErrInvalidMove
}

// Get all possible moves for the current player
func (g *Game) GetPossibleMoves() []Move {
	possibleMoves := []Move{}
	for _, row := range g.Board.Squares {
		for _, p := range row {
			if p != nil && p.Color == g.Turn {
				possibleMoves = append(possibleMoves, p.GetPossibleMoves(g)...)
			}
		}
	}
	possibleMoves = append(possibleMoves, g.getPossibleCastles()...)
	return possibleMoves
}

func (g *Game) getPossibleCastles() []Move {
	possibleMoves := []Move{}
	home_rank := 1
	if g.Turn == "black" {
		home_rank = 8
	}
	king, err := g.Board.GetPieceAtSquare('e', home_rank)
	if err != nil {
		return possibleMoves
	}
	// Kingside castle
	kingsideRook, err := g.Board.GetPieceAtSquare('h', home_rank)
	if err != nil {
		return possibleMoves
	}
	fSquare, err := g.Board.GetPieceAtSquare('f', home_rank)
	if err != nil {
		return possibleMoves
	}
	gSquare, err := g.Board.GetPieceAtSquare('g', home_rank)
	if err != nil {
		return possibleMoves
	}
	if kingsideRook != nil && kingsideRook.PieceType == Rook &&
		kingsideRook.Color == g.Turn && !kingsideRook.Moved &&
		king != nil && king.PieceType == King &&
		king.Color == g.Turn && !king.Moved &&
		fSquare == nil && gSquare == nil {
		possibleMoves = append(possibleMoves, Move{
			Castle: "short",
		})
	}
	// Queenside castle
	queensideRook, err := g.Board.GetPieceAtSquare('a', home_rank)
	if err != nil {
		return possibleMoves
	}
	bSquare, err := g.Board.GetPieceAtSquare('b', home_rank)
	if err != nil {
		return possibleMoves
	}
	cSquare, err := g.Board.GetPieceAtSquare('c', home_rank)
	if err != nil {
		return possibleMoves
	}
	dSquare, err := g.Board.GetPieceAtSquare('d', home_rank)
	if err != nil {
		return possibleMoves
	}
	if queensideRook != nil && queensideRook.PieceType == Rook &&
		queensideRook.Color == g.Turn && !queensideRook.Moved &&
		king != nil && king.PieceType == King &&
		king.Color == g.Turn && !king.Moved &&
		bSquare == nil && cSquare == nil && dSquare == nil {
		possibleMoves = append(possibleMoves, Move{
			Castle: "long",
		})
	}
	return possibleMoves
}

func (g *Game) changeTurn() {
	if g.Turn == "white" {
		g.Turn = "black"
	} else {
		g.Turn = "white"
	}
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

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
