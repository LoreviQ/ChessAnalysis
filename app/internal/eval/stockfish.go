package eval

import (
	"fmt"
	"strconv"
	"strings"
)

// change these for production
var (
	THREADS = "12"  // CPU Threads used by the engine
	HASH    = "256" // Size of hash table (MB)
	MultiPV = "1"   // Number of best lines found by the engine
	// SyzygyPath // ignore for now
)

// Sends the commands to set up stockfish 17 specifically returning the engine
func InitializeStockfish(filepath string, moveTime int) (*Engine, error) {
	eng, err := NewEngine(filepath, moveTime)
	if err != nil {
		return nil, err
	}
	eng.SendCommand("uci")
	for {
		response := eng.ReadResponse()
		if response[len(response)-1] == "uciok" {
			break
		}
	}
	// Set options
	eng.SendCommand(fmt.Sprintf("setoption name Threads value %v", THREADS))
	eng.SendCommand(fmt.Sprintf("setoption name Hash value %v", HASH))
	eng.SendCommand(fmt.Sprintf("setoption name MultiPV value %v", MultiPV))
	eng.SendCommand("isready")
	for {
		response := eng.ReadResponse()
		if response[len(response)-1] == "readyok" {
			break
		}
	}
	return eng, nil
}

// Evaluates the position string using the engine
//
// positionString is a space separated string of the moves in long algebraic notation
func (e *Engine) EvalPosition(positionString string) *MoveEval {
	e.SendCommand("ucinewgame")
	return e.queryPosition(positionString)
}

// Evaluates the game using the engine
//
// Returns an eval for each move in the game
func (e *Engine) EvalGame(positionString string) []*MoveEval {
	e.SendCommand("ucinewgame")
	moves := strings.Split(positionString, " ")
	eval := []*MoveEval{}
	for i := 0; i < len(moves)+1; i++ {
		eval = append(eval, e.queryPosition(strings.Join(moves[:i], " ")))
	}
	return eval
}

func (e *Engine) queryPosition(positionString string) *MoveEval {
	e.SendCommand(fmt.Sprintf("position startpos moves %v", positionString))
	e.SendCommand(fmt.Sprintf("go movetime %v", e.movetime))
	response := e.ReadResponse()
	eval, err := parseResponse(response)
	if err != nil {
		return nil
	}
	eval.moves = strings.Split(positionString, " ")
	return eval
}

// Parses the response from the engine
func parseResponse(response []string) (*MoveEval, error) {
	eval := &MoveEval{}
	// Penultimate line contains most recent info
	penultimateLine := strings.Split(response[len(response)-2], " ")
	for i, word := range penultimateLine {
		if word == "depth" {
			depth, err := strconv.Atoi(penultimateLine[i+1])
			if err != nil {
				return nil, err
			}
			eval.depth = depth
		}
		if word == "score" {
			if penultimateLine[i+1] == "mate" {
				eval.Mate = true
				mateIn, err := strconv.Atoi(penultimateLine[i+2])
				if err != nil {
					return nil, err
				}
				eval.MateIn = mateIn
			} else {
				score, err := strconv.Atoi(penultimateLine[i+2])
				if err != nil {
					return nil, err
				}
				eval.Score = score
			}
		}
		if word == "pv" {
			eval.bestLine = penultimateLine[i+1:]
		}
	}
	return eval, nil
}
