package eval

import (
	"bufio"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type Engine struct {
	cmd        *exec.Cmd
	writer     *bufio.Writer
	scanner    *bufio.Scanner
	Path       string
	Movetime   int    // ms spent on each move
	Depth      int    // max depth to search
	Threads    int    // number of threads to use
	Hash       int    // hash table size (MB)
	MultiPV    int    // number of lines to consider
	SyzygyPath string // path to syzygy tablebases
}

type MoveEval struct {
	Depth    int
	Score    int // centipawns
	BestLine []string
	Mate     bool
	MateIn   int
	PVnum    int
}

// NewEngine starts the provided engine and return a struct containing
// the command handle and input/output pipes
func NewEngine(filepath, Syzygy string, movetime, depth, threads, hash, multiPV int) (*Engine, error) {
	if filepath == "" {
		return nil, errors.New("no engine path provided")
	}
	cmd := exec.Command(filepath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(stdin)
	scanner := bufio.NewScanner(stdout)
	return &Engine{
		cmd:        cmd,
		writer:     writer,
		scanner:    scanner,
		Path:       filepath,
		Movetime:   movetime,
		Depth:      depth,
		Threads:    threads,
		Hash:       hash,
		MultiPV:    multiPV,
		SyzygyPath: Syzygy,
	}, nil
}

// SendCommand sends a command to the engine
func (e *Engine) SendCommand(command string) error {
	_, err := e.writer.WriteString(command + "\n")
	if err != nil {
		return err
	}
	return e.writer.Flush()
}

// ReadResponse reads the response from the engine
func (e *Engine) ReadResponse() []string {
	response := []string{}
	endStrings := []string{"uciok", "readyok", "bestmove", "Stockfish 17 by the Stockfish developers (see AUTHORS file)"}
	for e.scanner.Scan() {
		line := e.scanner.Text()
		response = append(response, line)
		for _, str := range endStrings {
			if strings.Contains(line, str) {
				return response
			}
		}
	}
	return response
}

// Close closes the engine
func (e *Engine) Close() error {
	if err := e.SendCommand("quit"); err != nil {
		return err
	}
	return e.cmd.Wait()
}

// Parses a score string and returns a MoveEval struct
//
// Expected input: "M#": mate in # moves
// Expected input: "#": centipawn score
func ParseScoreStr(scoreStr string) *MoveEval {
	if scoreStr[0] == 'M' {
		mateIn, _ := strconv.Atoi(scoreStr[1:])
		return &MoveEval{
			Mate:   true,
			MateIn: mateIn,
		}
	}
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return &MoveEval{}
	}
	return &MoveEval{
		Score: score,
		PVnum: 1,
	}
}

func GetEvalNum(evals []*MoveEval, pvNum int) *MoveEval {
	if len(evals) == 0 {
		return nil
	}
	for _, eval := range evals {
		if eval.PVnum == pvNum {
			return eval
		}
	}
	return evals[0]
}
