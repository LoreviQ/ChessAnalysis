package eval

import (
	"fmt"
	"strconv"
	"strings"
)

// Sends the commands to set up stockfish 17 specifically returning the engine
func InitializeStockfish(filepath, SyzygyPath string, moveTime, depth, threads, hash, MultiPV int) (*Engine, error) {
	if filepath == "" {
		return nil, fmt.Errorf("no engine path provided")
	}
	eng, err := NewEngine(filepath, SyzygyPath, moveTime, depth, threads, hash, MultiPV)
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
	eng.SendCommand(fmt.Sprintf("setoption name Threads value %d", threads))
	eng.SendCommand(fmt.Sprintf("setoption name Hash value %d", hash))
	eng.SendCommand(fmt.Sprintf("setoption name MultiPV value %d", MultiPV))
	if SyzygyPath != "" {
		eng.SendCommand(fmt.Sprintf("setoption name SyzygyPath value %v", SyzygyPath))
	}
	eng.SendCommand("setoption name UCI_ShowWDL value true")
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
func (e *Engine) EvalPosition(positionString string) []*MoveEval {
	e.SendCommand("ucinewgame")
	return e.queryPosition(positionString)
}

// Evaluates the game using the engine
//
// Returns an eval for each move in the game
func (e *Engine) EvalGame(positionString string) [][]*MoveEval {
	e.SendCommand("ucinewgame")
	moves := strings.Split(positionString, " ")
	gameEval := make([][]*MoveEval, len(moves)+1)
	for i := 0; i < len(moves)+1; i++ {
		evals := e.queryPosition(strings.Join(moves[:i], " "))
		gameEval[i] = evals
	}
	return gameEval
}

func (e *Engine) queryPosition(positionString string) []*MoveEval {
	e.SendCommand(fmt.Sprintf("position startpos moves %v", positionString))
	e.SendCommand(fmt.Sprintf("go depth %v movetime %v", e.Depth, e.Movetime))
	turnMult := 1
	if len(strings.Split(positionString, " "))%2 == 1 && positionString != "" {
		turnMult = -1
	}
	response := e.ReadResponse()
	evals, err := e.parseResponse(response, turnMult)
	if err != nil {
		return nil
	}
	return evals
}

// Parses the response from the engine
func (e *Engine) parseResponse(response []string, turnMult int) ([]*MoveEval, error) {
	evals := make([]*MoveEval, e.MultiPV)
	for i := range e.MultiPV {
		eval := &MoveEval{}
		if len(response) < 2+i {
			return nil, fmt.Errorf("not enough responses")
		}
		dataLine := strings.Split(response[len(response)-2-i], " ")
		for i, word := range dataLine {
			if word == "depth" {
				depth, err := strconv.Atoi(dataLine[i+1])
				if err != nil {
					return nil, err
				}
				eval.Depth = depth
			}
			if word == "score" {
				if dataLine[i+1] == "mate" {
					eval.Mate = true
					mateIn, err := strconv.Atoi(dataLine[i+2])
					if err != nil {
						return nil, err
					}
					eval.MateIn = mateIn
				} else {
					score, err := strconv.Atoi(dataLine[i+2])
					if err != nil {
						return nil, err
					}
					eval.Score = score * turnMult
				}
			}
			if word == "pv" {
				eval.BestLine = dataLine[i+1:]
			}
			if word == "multipv" {
				multiPV, err := strconv.Atoi(dataLine[i+1])
				if err != nil {
					return nil, err
				}
				eval.PVnum = multiPV
			}
		}
		evals[i] = eval
	}
	return evals, nil
}

func (e *Engine) ChangeOption(option, value string) error {
	switch option {
	case "MoveTime":
		moveTime, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		e.Movetime = moveTime
	case "Threads":
		threads, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		e.Threads = threads
	case "Hash":
		hash, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		e.Hash = hash
	case "MultiPV":
		multiPV, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		e.MultiPV = multiPV
	case "SyzygyPath":
		e.SyzygyPath = value
	}
	e.SendCommand(fmt.Sprintf("setoption name %v value %v", option, value))
	return nil
}
