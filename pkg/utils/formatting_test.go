package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "Empty string with zero width",
			input:    "",
			width:    0,
			expected: "",
		},
		{
			name:     "Empty string with positive width",
			input:    "",
			width:    5,
			expected: "     ",
		},
		{
			name:     "String shorter than width",
			input:    "abc",
			width:    5,
			expected: "abc  ",
		},
		{
			name:     "String equal to width",
			input:    "abcde",
			width:    5,
			expected: "abcde",
		},
		{
			name:     "String longer than width",
			input:    "abcdefg",
			width:    5,
			expected: "abcdefg",
		},
		{
			name:     "Negative width",
			input:    "abc",
			width:    -1,
			expected: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PadRight(tt.input, tt.width)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("PadRight(%q, %d) mismatch (-want +got):\n%s", tt.input, tt.width, diff)
			}
		})
	}
}

func TestSeparatorCell(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		justify  string
		expected string
	}{
		{
			name:     "Zero width",
			width:    0,
			justify:  "left",
			expected: "",
		},
		{
			name:     "Negative width",
			width:    -1,
			justify:  "left",
			expected: "",
		},
		{
			name:     "Left justify with width 1",
			width:    1,
			justify:  "left",
			expected: ":",
		},
		{
			name:     "Left justify with width 5",
			width:    5,
			justify:  "left",
			expected: ":----",
		},
		{
			name:     "Center justify with width 1",
			width:    1,
			justify:  "center",
			expected: ":-:",
		},
		{
			name:     "Center justify with width 2",
			width:    2,
			justify:  "center",
			expected: ":-:",
		},
		{
			name:     "Center justify with width 3",
			width:    3,
			justify:  "center",
			expected: ":-:",
		},
		{
			name:     "Center justify with width 5",
			width:    5,
			justify:  "center",
			expected: ":---:",
		},
		{
			name:     "Right justify with width 1",
			width:    1,
			justify:  "right",
			expected: ":",
		},
		{
			name:     "Right justify with width 5",
			width:    5,
			justify:  "right",
			expected: "----:",
		},
		{
			name:     "Invalid justify value",
			width:    5,
			justify:  "invalid",
			expected: ":----",
		},
		{
			name:     "Mixed case justify",
			width:    5,
			justify:  "CeNtEr",
			expected: ":---:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SeparatorCell(tt.width, tt.justify)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("SeparatorCell(%d, %q) mismatch (-want +got):\n%s", tt.width, tt.justify, diff)
			}
		})
	}
}
