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

// Converts a slice of moves to long algebraic notation
// by playing them and getting the long algebraic notation of the move history
func ConvertNotation(moves []string) ([]string, error) {
	g := NewGame()
	err := g.Moves(moves)
	if err != nil {
		return nil, err
	}
	return ConvertMovesToLongAlgebraicNotation(g.MoveHistory), nil
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
			fmt.Println("Type 'new_game' to start a new game")
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
		case "new_game":
			g.NewGame()
			continue
		default:
			_, err := g.Move(userInput)
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
func (g *Game) Move(moveStr string) (Move, error) {
	// Check if the move string is valid
	move, err := parseRegex(moveStr)
	if err != nil {
		return Move{}, err
	}
	// Get all possible moves for the current player
	possibleMoves := g.GetPossibleMoves()
	// Find the move that corresponds to the given move and play it
	correspondingMove, err := getCorrespondingMove(move, possibleMoves)
	if err != nil {
		return Move{}, err
	}
	// inherit check status from the move
	if move.CheckStatus != 0 {
		correspondingMove.CheckStatus = move.CheckStatus
	}
	if correspondingMove.Castle == "" {
		err = g.Board.MovePiece(correspondingMove)
	} else {
		err = g.Castle(correspondingMove.Castle)
	}
	if err != nil {
		return Move{}, err
	}
	g.MoveHistory = append(g.MoveHistory, correspondingMove)
	g.changeTurn()
	return correspondingMove, nil
}

// Takes a slice of move strings in algebraic notation and plays them
func (g *Game) Moves(moveStrs []string) error {
	for _, moveStr := range moveStrs {
		_, err := g.Move(moveStr)
		if err != nil {
			return err
		}
	}
	return nil
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
	homeRank := 1
	if g.Turn == "black" {
		homeRank = 8
	}
	king, err := g.Board.GetPieceAtSquare('e', homeRank)
	if err != nil {
		return possibleMoves
	}
	// Kingside castle
	kingsideRook, err := g.Board.GetPieceAtSquare('h', homeRank)
	if err != nil {
		return possibleMoves
	}
	fSquare, err := g.Board.GetPieceAtSquare('f', homeRank)
	if err != nil {
		return possibleMoves
	}
	gSquare, err := g.Board.GetPieceAtSquare('g', homeRank)
	if err != nil {
		return possibleMoves
	}
	if kingsideRook != nil && kingsideRook.PieceType == Rook &&
		kingsideRook.Color == g.Turn && !kingsideRook.Moved &&
		king != nil && king.PieceType == King &&
		king.Color == g.Turn && !king.Moved &&
		fSquare == nil && gSquare == nil {
		possibleMoves = append(possibleMoves, Move{
			FromRank: homeRank,
			Castle:   "short",
		})
	}
	// Queenside castle
	queensideRook, err := g.Board.GetPieceAtSquare('a', homeRank)
	if err != nil {
		return possibleMoves
	}
	bSquare, err := g.Board.GetPieceAtSquare('b', homeRank)
	if err != nil {
		return possibleMoves
	}
	cSquare, err := g.Board.GetPieceAtSquare('c', homeRank)
	if err != nil {
		return possibleMoves
	}
	dSquare, err := g.Board.GetPieceAtSquare('d', homeRank)
	if err != nil {
		return possibleMoves
	}
	if queensideRook != nil && queensideRook.PieceType == Rook &&
		queensideRook.Color == g.Turn && !queensideRook.Moved &&
		king != nil && king.PieceType == King &&
		king.Color == g.Turn && !king.Moved &&
		bSquare == nil && cSquare == nil && dSquare == nil {
		possibleMoves = append(possibleMoves, Move{
			FromRank: homeRank,
			Castle:   "long",
		})
	}
	return possibleMoves
}

func (g *Game) Castle(castletype string) error {
	homeRank := 1
	var err error
	b := g.Board
	king, err := b.GetPieceAtSquare('e', homeRank)
	if err != nil {
		return err
	}
	if g.Turn == "black" {
		homeRank = 8
	}
	if castletype == "short" {
		err = b.MovePiece(Move{
			Piece:    'K',
			FromFile: 'e',
			FromRank: homeRank,
			ToFile:   'g',
			ToRank:   homeRank,
		})
		if err != nil {
			return err
		}
		err = b.MovePiece(Move{
			Piece:    'R',
			FromFile: 'h',
			FromRank: homeRank,
			ToFile:   'f',
			ToRank:   homeRank,
		})
		if err != nil {
			// Undo the king move
			b.MovePiece(Move{
				Piece:    'K',
				FromFile: 'g',
				FromRank: homeRank,
				ToFile:   'e',
				ToRank:   homeRank,
			})
			king.Moved = false
			return err
		}
	} else if castletype == "long" {
		err = b.MovePiece(Move{
			Piece:    'K',
			FromFile: 'e',
			FromRank: homeRank,
			ToFile:   'c',
			ToRank:   homeRank,
		})
		if err != nil {
			return err
		}
		err = b.MovePiece(Move{
			Piece:    'R',
			FromFile: 'a',
			FromRank: homeRank,
			ToFile:   'd',
			ToRank:   homeRank,
		})
		if err != nil {
			// Undo the king move
			b.MovePiece(Move{
				Piece:    'K',
				FromFile: 'c',
				FromRank: homeRank,
				ToFile:   'e',
				ToRank:   homeRank,
			})
			king.Moved = false
			return err
		}
	} else {
		return ErrInvalidMove
	}
	return nil
}

func (g *Game) changeTurn() {
	if g.Turn == "white" {
		g.Turn = "black"
	} else {
		g.Turn = "white"
	}
}

func (g *Game) NewGame() {
	g.Board = NewBoard()
	g.Turn = "white"
	g.MoveHistory = []Move{}
}

func (g *Game) Clone() *Game {
	return &Game{
		Board:       g.Board.Clone(),
		Turn:        g.Turn,
		MoveHistory: g.MoveHistory,
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
