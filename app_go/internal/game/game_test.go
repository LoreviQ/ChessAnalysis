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
