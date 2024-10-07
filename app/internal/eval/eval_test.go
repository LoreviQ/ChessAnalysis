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

func TestInitializeEngine(t *testing.T) {
	eng, err := InitializeStockfish()
	if err != nil {
		t.Errorf("InitializeStockfish() failed: %v", err)
	}
	if eng == nil {
		t.Errorf("InitializeStockfish() failed: returned nil")
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

func TestEvalPosition(t *testing.T) {
	eng, err := InitializeStockfish()
	if err != nil {
		t.Errorf("InitializeStockfish() failed: %v", err)
	}
	eval := eng.EvalPosition("e2e4 e7e5 b1c3 b8c6 f2f4 e5f4 g1f3 f8b4 d2d4 b4c3 b2c3 d7d5 e4e5 f7f6 c1f4")
	if eval == nil {
		t.Errorf("EvalPosition() failed: returned nil")
	} else {
		if eval.depth == 0 {
			t.Errorf("EvalPosition() failed: expected depth != 0, got %v", eval.depth)
		}
		if eval.score == 0 {
			t.Errorf("EvalPosition() failed: expected score != 0, got %v", eval.score)
		}
		if eval.bestLine == nil {
			t.Error("EvalPosition() failed: expected bestLine != nil")
		}
	}
}
