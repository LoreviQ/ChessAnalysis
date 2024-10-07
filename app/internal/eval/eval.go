package eval

import (
	"bufio"
	"os/exec"
	"strings"
)

type Engine struct {
	cmd     *exec.Cmd
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

// StartStockfish starts the Stockfish engine and returns the command handle and input/output pipes
func StartStockfish(filepath string) (*Engine, error) {
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
	return &Engine{cmd, writer, scanner}, nil
}

// SendCommand sends a command to the Stockfish engine
func (e *Engine) SendCommand(command string) error {
	_, err := e.writer.WriteString(command + "\n")
	if err != nil {
		return err
	}
	return e.writer.Flush()
}

// ReadResponse reads the response from the Stockfish engine
func (e *Engine) ReadResponse() string {
	response := ""
	for e.scanner.Scan() {
		line := e.scanner.Text()
		response += line + "\n"
		if line == "uciok" || line == "readyok" || strings.Contains(line, "bestmove") {
			break
		}
	}
	return response
}
