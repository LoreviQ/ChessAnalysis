package eval

import (
	"strings"
	"testing"
)

var (
	FILEPATH   = "/home/lorevi/workspace/stockfish/stockfish-ubuntu-x86-64-avx2"
	SYZYGYPATH = "/home/lorevi/workspace/3-4-5"
	MOVETIME   = 60
	THREADS    = 12
	HASH       = 256
	MULTIPV    = 1
)

func TestNewEngine(t *testing.T) {
	eng, err := NewEngine(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, MULTIPV)
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
	eng, err := InitializeStockfish(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, MULTIPV)
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
	eng, err := NewEngine(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, MULTIPV)
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
	eng, err := InitializeStockfish(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, MULTIPV)
	if err != nil {
		t.Errorf("InitializeStockfish() failed: %v", err)
	}
	eval := eng.EvalPosition("e2e4 e7e5 b1c3 b8c6 f2f4 e5f4 g1f3 f8b4 d2d4 b4c3 b2c3 d7d5 e4e5 f7f6 c1f4")
	if eval == nil {
		t.Errorf("EvalPosition() failed: returned nil")
	} else {
		if eval[0].Depth == 0 {
			t.Errorf("EvalPosition() failed: expected depth != 0, got %v", eval[0].Depth)
		}
		if eval[0].Score == 0 {
			t.Errorf("EvalPosition() failed: expected score != 0, got %v", eval[0].Score)
		}
		if eval[0].BestLine == nil {
			t.Error("EvalPosition() failed: expected bestLine != nil")
		}
	}
}

func TestEvalGame(t *testing.T) {
	eng, err := InitializeStockfish(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, MULTIPV)
	if err != nil {
		t.Errorf("InitializeStockfish() failed: %v", err)
	}
	positionString := "e2e4 e7e5 b1c3 b8c6 f2f4 e5f4 g1f3 f8b4 d2d4 b4c3 b2c3 d7d5 e4e5 f7f6 c1f4"
	eval := eng.EvalGame(positionString)
	expected := strings.Split(positionString, " ")
	if len(eval) != len(expected)+1 {
		t.Errorf("EvalGame() failed: expected %v moves, got %v", len(expected)+1, len(eval))
	}
	for i, moveEval := range eval {
		if moveEval == nil {
			t.Errorf("EvalGame() failed: expected moveEval != nil at index %v", i)
		} else {
			if moveEval[0].Depth == 0 {
				t.Errorf("EvalPosition() failed: expected depth != 0, got %v", moveEval[0].Depth)
			}
			if moveEval[0].Score == 0 {
				t.Errorf("EvalPosition() failed: expected score != 0, got %v", moveEval[0].Score)
			}
			if moveEval[0].BestLine == nil {
				t.Error("EvalPosition() failed: expected bestLine != nil")
			}
		}
	}
}

func TestMultiPV(t *testing.T) {
	eng, err := InitializeStockfish(FILEPATH, SYZYGYPATH, MOVETIME, THREADS, HASH, 3)
	if err != nil {
		t.Errorf("InitializeStockfish() failed: %v", err)
	}
	positionString := "e2e4 e7e5 b1c3 b8c6 f2f4 e5f4 g1f3 f8b4 d2d4 b4c3 b2c3 d7d5 e4e5 f7f6 c1f4"
	eval := eng.EvalGame(positionString)
	for i, moveEval := range eval {
		if len(moveEval) != 3 {
			t.Errorf("EvalGame() failed: expected 3 moveEvals, got %v at index %v", len(moveEval), i)
		}
		for i, mPV := range moveEval {
			if mPV == nil {
				t.Errorf("EvalGame() failed: expected moveEval != nil at index %v", i)
			} else {
				if mPV.PVnum != 1 && mPV.PVnum != 2 && mPV.PVnum != 3 {
					t.Errorf("EvalGame() failed: expected PVnum to be 1, 2, or 3, got %v", mPV.PVnum)
				}
			}
		}
	}
}
