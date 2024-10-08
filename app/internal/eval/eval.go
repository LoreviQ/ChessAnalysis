package eval

import (
	"bufio"
	"os/exec"
	"strings"
)

type Engine struct {
	cmd      *exec.Cmd
	writer   *bufio.Writer
	scanner  *bufio.Scanner
	movetime int // ms spent on each move
}

type MoveEval struct {
	moves    []string
	Depth    int
	Score    int // centipawns
	bestLine []string
	Mate     bool
	MateIn   int
}

// NewEngine starts the provided engine and return a struct containing
// the command handle and input/output pipes
func NewEngine(filepath string, movetime int) (*Engine, error) {
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
		cmd:      cmd,
		writer:   writer,
		scanner:  scanner,
		movetime: movetime,
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
