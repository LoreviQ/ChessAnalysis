package game

import (
	"testing"
)

func TestLongAlgebraicNotation(t *testing.T) {
	tests := []struct {
		move     Move
		expected string
		err      error
	}{
		{
			move: Move{
				Piece:    'R',
				FromFile: 'g',
				FromRank: 4,
				ToFile:   'e',
				ToRank:   4,
			},
			expected: "Rg4e4",
			err:      nil,
		},
		{
			move: Move{
				FromFile: 'e',
				FromRank: 2,
				ToFile:   'e',
				ToRank:   4,
			},
			expected: "e2e4",
			err:      nil,
		},
		{
			move: Move{
				Piece:    'N',
				FromFile: 'g',
				FromRank: 1,
				ToFile:   'f',
				ToRank:   3,
			},
			expected: "Ng1f3",
			err:      nil,
		},
		{
			move: Move{
				FromFile: 'e',
				FromRank: 4,
				Capture:  'x',
				ToFile:   'd',
				ToRank:   5,
			},
			expected: "e4xd5",
			err:      nil,
		},
		{
			move: Move{
				FromFile:  'e',
				FromRank:  7,
				ToFile:    'e',
				ToRank:    8,
				Promotion: 'Q',
			},
			expected: "e7e8=Q",
			err:      nil,
		},
		{
			move: Move{
				Castle: "long",
			},
			expected: "O-O-O",
			err:      nil,
		},
		{
			move: Move{
				Castle: "short",
			},
			expected: "O-O",
			err:      nil,
		},
		// Failing cases
		{
			move: Move{
				ToFile: 'e',
				ToRank: 4,
			},
			expected: "",
			err:      ErrNotEnoughInfo,
		},
		{
			move: Move{
				Piece:  'R',
				ToFile: 'g',
				ToRank: 5,
			},
			expected: "",
			err:      ErrNotEnoughInfo,
		},
	}

	for _, tt := range tests {
		actual, err := tt.move.LongAlgebraicNotation()
		if err != tt.err {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != tt.expected {
			t.Errorf("Expected: %s, got: %s", tt.expected, actual)
		}
	}
}
func TestShortAlgebraicNotation(t *testing.T) {
	tests := []struct {
		move     Move
		expected string
		err      error
	}{
		{
			move: Move{
				Piece:    'R',
				FromFile: 'g',
				FromRank: 4,
				ToFile:   'e',
				ToRank:   4,
			},
			expected: "Re4",
			err:      nil,
		},
		{
			move: Move{
				FromFile: 'e',
				FromRank: 2,
				ToFile:   'e',
				ToRank:   4,
			},
			expected: "e4",
			err:      nil,
		},
		{
			move: Move{
				Piece:    'N',
				FromFile: 'g',
				FromRank: 1,
				ToFile:   'f',
				ToRank:   3,
			},
			expected: "Nf3",
			err:      nil,
		},
		{
			move: Move{
				FromFile: 'e',
				FromRank: 4,
				Capture:  'x',
				ToFile:   'd',
				ToRank:   5,
			},
			expected: "exd5",
			err:      nil,
		},
		{
			move: Move{
				FromFile:  'e',
				FromRank:  7,
				ToFile:    'e',
				ToRank:    8,
				Promotion: 'Q',
			},
			expected: "e8=Q",
			err:      nil,
		},
		{
			move: Move{
				Castle: "long",
			},
			expected: "O-O-O",
			err:      nil,
		},
		{
			move: Move{
				Castle: "short",
			},
			expected: "O-O",
			err:      nil,
		},
		// Failing cases
		{
			move: Move{
				ToFile: 'e',
			},
			expected: "",
			err:      ErrNotEnoughInfo,
		},
		{
			move: Move{
				Piece:  'R',
				ToRank: 5,
			},
			expected: "",
			err:      ErrNotEnoughInfo,
		},
	}

	for _, tt := range tests {
		actual, err := tt.move.ShortAlgebraicNotation(false, false)
		if err != tt.err {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != tt.expected {
			t.Errorf("Expected: %s, got: %s", tt.expected, actual)
		}
	}
}

func TestConvertMovesToShortAlgebraicNotation(t *testing.T) {
	moves := []Move{
		{
			Piece:    'R',
			FromFile: 'a',
			FromRank: 3,
			ToFile:   'c',
			ToRank:   3,
		},
		{
			Piece:    'R',
			FromFile: 'h',
			FromRank: 3,
			ToFile:   'c',
			ToRank:   3,
		},
		{
			FromFile: 'c',
			FromRank: 2,
			ToFile:   'c',
			ToRank:   3,
		},
	}
	expected := []string{"Rac3", "Rhc3", "c3"}
	actual, err := ConvertMovesToShortAlgebraicNotation(moves)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("Expected %d moves, got %d", len(expected), len(actual))
	}
	for _, expectedString := range expected {
		found := false
		for actualString := range actual {
			if expectedString == actualString {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %s not found", expectedString)
		}
	}

}
