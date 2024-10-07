package eval

import (
	"testing"
)

var FILEPATH = "/home/lorevi/workspace/github.com/LoreviQ/stockfish/stockfish-ubuntu-x86-64-avx2"

func TestNewEngine(t *testing.T) {
	eng, err := NewEngine(FILEPATH)
	if err != nil {
		t.Errorf("NewEngine(stockfish) failed: %v", err)
	}
	if eng == nil {
		t.Errorf("NewEngine(stockfish) failed: returned nil")
	}
	err = eng.Close()
	if err != nil {
		t.Errorf("Engine.Close() failed: %v", err)
	}
}

// Tests both SendCommand and ReadResponse checking received response to see if
// it matches the expected response from sending the command
func TestSendCommandReadResponse(t *testing.T) {
	eng, err := NewEngine(FILEPATH)
	if err != nil {
		t.Errorf("NewEngine(stockfish) failed: %v", err)
	}
	// initial message from stockfish
	response := eng.ReadResponse()
	expected := "Stockfish 17 by the Stockfish developers (see AUTHORS file)"
	if response[len(response)-1] != expected {
		t.Errorf("ReadResponse() failed: expected %v, got %v", expected, response)
	}

	// send uci command
	err = eng.SendCommand("uci")
	if err != nil {
		t.Errorf("SendCommand(uci) failed: %v", err)
	}
	response = eng.ReadResponse()
	expected = "uciok"
	if response[len(response)-1] != expected {
		t.Errorf("ReadResponse() failed: expected %v, got %v", expected, response)
	}

	// send isready command
	err = eng.SendCommand("isready")
	if err != nil {
		t.Errorf("SendCommand(isready) failed: %v", err)
	}
	response = eng.ReadResponse()
	expected = "readyok"
	if response[len(response)-1] != expected {
		t.Errorf("ReadResponse() failed: expected %v, got %v", expected, response)
	}

	// send commands for moves
	eng.SendCommand("ucinewgame")
	eng.SendCommand("position startpos moves e2e4 e7e5")
	eng.SendCommand("go movetime 1000ms")
	response = eng.ReadResponse()
	expected = "bestmove g1f3 ponder b8c6"
	if response[len(response)-1] != expected {
		t.Errorf("ReadResponse() failed: expected %v, got %v", expected, response[len(response)-1])
	}
}
