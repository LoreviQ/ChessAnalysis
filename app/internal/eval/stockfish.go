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

type moveEval struct {
	moves    []string
	depth    int
	score    int // centipawns
	bestLine []string
}

// Sends the commands to set up stockfish 17 specifically returning the engine
func InitializeStockfish() (*Engine, error) {
	eng, err := NewEngine("/home/lorevi/workspace/github.com/LoreviQ/stockfish/stockfish-ubuntu-x86-64-avx2")
	if err != nil {
		return nil, err
	}
	eng.SendCommand("uci")
	eng.ReadResponse()
	// Set options
	eng.SendCommand(fmt.Sprintf("setoption name Threads value %v", THREADS))
	eng.SendCommand(fmt.Sprintf("setoption name Hash value %v", HASH))
	eng.SendCommand(fmt.Sprintf("setoption name MultiPV value %v", MultiPV))
	eng.SendCommand("isready")
	eng.ReadResponse()
	return eng, nil
}

// Evaluates the position string using the engine
//
// positionString is a space separated string of the moves in long algebraic notation
func (e *Engine) EvalPosition(positionString string) *moveEval {
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
func parseResponse(response []string) (*moveEval, error) {
	eval := &moveEval{}
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
			score, err := strconv.Atoi(penultimateLine[i+2])
			if err != nil {
				return nil, err
			}
			eval.score = score
		}
		if word == "pv" {
			eval.bestLine = penultimateLine[i+1:]
		}
	}
	return eval, nil
}
