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
		for _, notation := range notations {
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
		{"e4", nil},
		{"e5", nil},
		{"Nc3", nil},
		{"Nf6", nil},
		{"f4", nil},
		{"exf4", nil},
		{"e5", nil},
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
