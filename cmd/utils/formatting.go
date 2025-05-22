package utils

import (
	"strings"
)

// PadRight pads a string with spaces to the specified width
func PadRight(s string, width int) string {
	if len(s) < width {
		return s + strings.Repeat(" ", width-len(s))
	}
	return s
}

// SeparatorCell creates a separator cell for markdown tables based on justification
func SeparatorCell(width int, justify string) string {
	if width <= 0 {
		return ""
	}
	switch strings.ToLower(justify) {
	case "center":
		if width <= 3 {
			return ":-:"
		}
		return ":" + strings.Repeat("-", width-2) + ":"
	case "right":
		return strings.Repeat("-", width-1) + ":"
	default: // left
		return ":" + strings.Repeat("-", width-1)
	}
}
