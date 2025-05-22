package calendar

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewOptions(t *testing.T) {
	// Test that NewOptions returns the expected default values
	options := NewOptions()

	// Create expected options with current year
	expected := Options{
		Year:             time.Now().Year(),
		Month:            nil,
		EndYear:          nil,
		EndMonth:         nil,
		FirstDayOfWeek:   time.Monday,
		ShowCalendarWeek: true,
		ShowWeekends:     true,
		ShowComments:     true,
		UseShortDayNames: false,
		Justify:          "left",
	}

	// Compare using cmp.Diff
	if diff := cmp.Diff(expected, options); diff != "" {
		t.Errorf("NewOptions() mismatch (-want +got):\n%s", diff)
	}
}
