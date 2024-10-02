package game

import (
	"reflect"
	"testing"
)

func TestParseRegex(t *testing.T) {
	tests := []struct {
		input    string
		expected Move
	}{
		{
			input: "Rgxe4",
			expected: Move{
				Piece:       'R',
				FromFile:    'g',
				FromRank:    0,
				Capture:     'x',
				ToFile:      'e',
				ToRank:      4,
				Promotion:   0,
				CheckStatus: 0,
				Castle:      "",
			},
		},
		{
			input: "O-O-O",
			expected: Move{
				Castle: "long",
			},
		},
		{
			input: "e4",
			expected: Move{
				ToFile: 'e',
				ToRank: 4,
			},
		},
		{
			input: "Nf3",
			expected: Move{
				Piece:  'N',
				ToFile: 'f',
				ToRank: 3,
			},
		},
		{
			input: "exd5",
			expected: Move{
				FromFile: 'e',
				Capture:  'x',
				ToFile:   'd',
				ToRank:   5,
			},
		},
		{
			input: "e8=Q",
			expected: Move{
				ToFile:    'e',
				ToRank:    8,
				Promotion: 'Q',
			},
		},
		{
			input: "O-O",
			expected: Move{
				Castle: "short",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseRegex(test.input)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestPossibleMoves(t *testing.T) {
	g := NewGame()
	moves := g.GetPossibleMoves()
	expectedMoves := []string{
		"a3", "a4", "b3", "b4",
		"c3", "c4", "d3", "d4",
		"e3", "e4", "f3", "f4",
		"g3", "g4", "h3", "h4",
		"Nc3", "Nf3", "Nh3", "Na3",
	}
	if len(moves) != 20 {
		t.Errorf("Expected 20 possible moves, got %d", len(moves))
	}
	notations, err := ConvertMovesToShortAlgebraicNotation(moves)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	for _, move := range expectedMoves {
		found := false
		for notation := range notations {
			if move == notation {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %s not found", move)
		}
	}
}

func TestPossibleMovesDuplicate(t *testing.T) {
	g := NewGame()
	b := g.Board
	// Move rooks up to test multiple pieces of the same type moving to the same square
	b.MovePiece(Move{FromFile: 'a', FromRank: 1, ToFile: 'a', ToRank: 3})
	b.MovePiece(Move{FromFile: 'h', FromRank: 1, ToFile: 'h', ToRank: 3})
	moves := g.GetPossibleMoves()
	expectedMoves := []string{
		// Pawn moves
		"b3", "b4", "c3", "c4",
		"d3", "d4", "e3", "e4",
		"f3", "f4", "g3", "g4",
		// Knight moves
		"Nc3", "Nf3",
		// Basic rook moves
		"Ra4", "Rh4", "Ra5", "Rh5",
		"Ra6", "Rh6", "Rxa7", "Rxh7",
		// Rook moves to the same square
		"Rab3", "Rac3", "Rad3", "Rae3",
		"Raf3", "Rag3", "Rhb3", "Rhc3",
		"Rhd3", "Rhe3", "Rhf3", "Rhg3",
	}
	if len(moves) != len(expectedMoves) {
		t.Errorf("Expected 32 possible moves, got %d", len(moves))
	}
	notations, err := ConvertMovesToShortAlgebraicNotation(moves)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	for _, move := range expectedMoves {
		found := false
		for notation := range notations {
			if move == notation {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %s not found", move)
		}
	}
}

func TestPossibleMovesCastle(t *testing.T) {
	g := NewGame()
	b := g.Board
	// Move pieces up to open up castling
	b.MovePiece(Move{FromFile: 'b', FromRank: 1, ToFile: 'b', ToRank: 3})
	b.MovePiece(Move{FromFile: 'c', FromRank: 1, ToFile: 'c', ToRank: 3})
	b.MovePiece(Move{FromFile: 'd', FromRank: 1, ToFile: 'd', ToRank: 3})
	b.MovePiece(Move{FromFile: 'f', FromRank: 1, ToFile: 'f', ToRank: 3})
	b.MovePiece(Move{FromFile: 'g', FromRank: 1, ToFile: 'g', ToRank: 3})
	moves := g.GetPossibleMoves()
	expectedMoves := []string{
		// Pawn moves
		"a3", "a4", "e3", "e4", "h3", "h4",
		// Knight moves
		"Nc1", "Na5", "Nc5", "Nd4",
		"Nf1", "Nf5", "Nh5", "Ne4",
		// Bishop moves
		"Bb4", "Ba5", "Bd4", "Be5", "Bf6", "Bxg7",
		"Bg4", "Bh5", "Be4", "Bd5", "Bc6", "Bxb7",
		// Rook moves
		"Rb1", "Rc1", "Rd1", "Rf1", "Rg1",
		// Queen moves
		"Qd4", "Qd5", "Qd6", "Qxd7",
		"Qe4", "Qf5", "Qg6", "Qxh7",
		"Qc4", "Qb5", "Qa6", "Qe3",
		// King moves
		"Kd1", "Kf1",
		// Castling
		"O-O", "O-O-O",
	}
	if len(moves) != len(expectedMoves) {
		t.Errorf("Expected %d possible moves, got %d", len(expectedMoves), len(moves))
	}
	notations, err := ConvertMovesToShortAlgebraicNotation(moves)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	for _, move := range expectedMoves {
		found := false
		for notation := range notations {
			if move == notation {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %s not found", move)
		}
	}
}

func TestMove(t *testing.T) {
	g := NewGame()
	tests := []struct {
		move string
		err  error
	}{
		// First 7 moves of Vienna Game followed by 2 invalid moves
		{"e4", nil},
		{"e5", nil},
		{"Nc3", nil},
		{"Nf6", nil},
		{"f4", nil},
		{"exf4", nil},
		{"e5", nil},
		{"Ne1", ErrInvalidMove},
		{"Qh4", ErrInvalidMove},
	}

	for _, tt := range tests {
		t.Run(tt.move, func(t *testing.T) {
			err := g.Move(tt.move)
			if err != tt.err {
				t.Errorf("Expected: %v, got: %v", tt.err, err)
			}
		})
	}
}

func TestMoveAmbiguous(t *testing.T) {
	g := NewGame()
	tests := []struct {
		move string
		err  error
	}{
		// Quick move order to produce an ambiguous move
		{"a4", nil},
		{"a5", nil},
		{"h4", nil},
		{"h5", nil},
		{"Ra3", nil},
		{"Ra6", nil},
		{"Rh3", ErrAmbiguousMove},
		{"Rhh3", nil},
	}

	for _, tt := range tests {
		t.Run(tt.move, func(t *testing.T) {
			err := g.Move(tt.move)
			if err != tt.err {
				t.Errorf("Expected: %v, got: %v", tt.err, err)
			}
		})
	}
}

func TestConvertNotation(t *testing.T) {
	moves_to_convert := []string{
		"e4", "e5", "Nc3", "Nc6", "f4", "exf4", "Nf3", "Bb4",
		"d4", "Bxc3+", "bxc3", "d5", "e5", "f6", "Bxf4", "fxe5",
		"Bxe5", "Nxe5", "Nxe5", "Qe7", "Bd3", "c5", "O-O", "Nf6",
		"Qf3", "Bg4", "Qf2", "O-O", "Qg3", "Ne4", "Bxe4", "dxe4",
		"Qxg4", "cxd4", "cxd4", "Rad8", "c3", "b5", "Qxe4", "Qa3",
		"Rf3", "a5", "Raf1", "Qxa2", "Rxf8+", "Rxf8", "Rxf8+",
		"Kxf8", "Qa8+", "Ke7", "Nc6+", "Ke6", "Nxa5", "Qe2",
		"Qc6+", "Kf5", "Qf3+", "Qxf3", "gxf3", "Kf4", "Kf2",
		"g5", "d5", "Ke5", "Nc6+", "Kxd5", "Nd4", "b4", "cxb4",
		"Kxd4", "Kg3", "Kc4", "Kg4", "h6", "Kh5", "Kxb4",
		"Kxh6", "Kc5", "Kxg5", "Kd6", "Kg6", "Ke5", "h4",
		"Kf4", "h5", "Kxf3", "h6", "Ke2", "h7", "Kd3", "h8=Q",
		"Kc4", "Qe5", "Kd3", "Kf5", "Kc4", "Kf4", "Kd3",
		"Qe4+", "Kc3", "Kf3", "Kd2", "Qe3+", "Kc2", "Kf2",
		"Kd1", "Qe2+", "Kc1", "Ke3", "Kb1", "Kd3", "Kc1",
		"Qf2", "Kd1", "Qf1#",
	}
	expected_moves := []string{
		"e2e4", "e7e5", "Nb1c3", "Nb8c6", "f2f4", "e5xf4", "Ng1f3", "Bf8b4",
		"d2d4", "Bb4xc3+", "b2xc3", "d7d5", "e4e5", "f7f6", "Bc1xf4", "f6xe5",
		"Bf4xe5", "Nc6xe5", "Nf3xe5", "Qd8e7", "Bf1d3", "c7c5", "O-O", "Ng8f6",
		"Qd1f3", "Bc8g4", "Qf3f2", "O-O", "Qf2g3", "Nf6e4", "Bd3xe4", "d5xe4",
		"Qg3xg4", "c5xd4", "c3xd4", "Ra8d8", "c2c3", "b7b5", "Qg4xe4", "Qe7a3",
		"Rf1f3", "a7a5", "Ra1f1", "Qa3xa2", "Rf3xf8+", "Rd8xf8", "Rf1xf8+",
		"Kg8xf8", "Qe4a8+", "Kf8e7", "Ne5c6+", "Ke7e6", "Nc6xa5", "Qa2e2",
		"Qa8c6+", "Ke6f5", "Qc6f3+", "Qe2xf3", "g2xf3", "Kf5f4", "Kg1f2",
		"g7g5", "d4d5", "Kf4e5", "Na5c6+", "Ke5xd5", "Nc6d4", "b5b4", "c3xb4",
		"Kd5xd4", "Kf2g3", "Kd4c4", "Kg3g4", "h7h6", "Kg4h5", "Kc4xb4",
		"Kh5xh6", "Kb4c5", "Kh6xg5", "Kc5d6", "Kg5g6", "Kd6e5", "h2h4",
		"Ke5f4", "h4h5", "Kf4xf3", "h5h6", "Kf3e2", "h6h7", "Ke2d3", "h7h8=Q",
		"Kd3c4", "Qh8e5", "Kc4d3", "Kg6f5", "Kd3c4", "Kf5f4", "Kc4d3",
		"Qe5e4+", "Kd3c3", "Kf4f3", "Kc3d2", "Qe4e3+", "Kd2c2", "Kf3f2",
		"Kc2d1", "Qe3e2+", "Kd1c1", "Kf2e3", "Kc1b1", "Ke3d3", "Kb1c1",
		"Qe2f2", "Kc1d1", "Qf2f1#",
	}
	converted_moves, err := ConvertNotation(moves_to_convert)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(converted_moves) != len(expected_moves) {
		t.Errorf("Expected %d moves, got %d", len(expected_moves), len(converted_moves))
	}
	for i, move := range expected_moves {
		if converted_moves[i] != move {
			t.Errorf("Expected %s, got %s", move, converted_moves[i])
		}
	}
}
