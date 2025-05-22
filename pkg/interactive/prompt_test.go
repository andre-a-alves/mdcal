package interactive

import (
	"testing"
	"time"

	"github.com/andre-a-alves/mdcal/pkg/calendar"
)

func TestInitialize(t *testing.T) {
	// Test with default options
	options := calendar.NewOptions()
	model := Initialize(&options)

	// Check that the model was initialized correctly
	if model.weekdayIndex != int(options.FirstDayOfWeek) {
		t.Errorf("Expected weekdayIndex to be %d, got %d", int(options.FirstDayOfWeek), model.weekdayIndex)
	}

	if model.showWeekNumbers != options.ShowCalendarWeek {
		t.Errorf("Expected showWeekNumbers to be %v, got %v", options.ShowCalendarWeek, model.showWeekNumbers)
	}

	if model.showWeekends != options.ShowWeekends {
		t.Errorf("Expected showWeekends to be %v, got %v", options.ShowWeekends, model.showWeekends)
	}

	if model.showComments != options.ShowComments {
		t.Errorf("Expected showComments to be %v, got %v", options.ShowComments, model.showComments)
	}

	if model.useShortDayNames != options.UseShortDayNames {
		t.Errorf("Expected useShortDayNames to be %v, got %v", options.UseShortDayNames, model.useShortDayNames)
	}

	if model.dateRangeEnabled != (options.EndMonth != nil) {
		t.Errorf("Expected dateRangeEnabled to be %v, got %v", options.EndMonth != nil, model.dateRangeEnabled)
	}

	// Test with custom options
	month := 5
	endYear := 2024
	endMonth := 7
	customOptions := calendar.Options{
		Year:             2023,
		Month:            &month,
		EndYear:          &endYear,
		EndMonth:         &endMonth,
		FirstDayOfWeek:   time.Sunday,
		ShowCalendarWeek: false,
		ShowWeekends:     false,
		ShowComments:     false,
		UseShortDayNames: true,
		Justify:          "center",
	}
	customModel := Initialize(&customOptions)

	// Check that the model was initialized correctly
	if customModel.weekdayIndex != int(customOptions.FirstDayOfWeek) {
		t.Errorf("Expected weekdayIndex to be %d, got %d", int(customOptions.FirstDayOfWeek), customModel.weekdayIndex)
	}

	if customModel.showWeekNumbers != customOptions.ShowCalendarWeek {
		t.Errorf("Expected showWeekNumbers to be %v, got %v", customOptions.ShowCalendarWeek, customModel.showWeekNumbers)
	}

	if customModel.showWeekends != customOptions.ShowWeekends {
		t.Errorf("Expected showWeekends to be %v, got %v", customOptions.ShowWeekends, customModel.showWeekends)
	}

	if customModel.showComments != customOptions.ShowComments {
		t.Errorf("Expected showComments to be %v, got %v", customOptions.ShowComments, customModel.showComments)
	}

	if customModel.useShortDayNames != customOptions.UseShortDayNames {
		t.Errorf("Expected useShortDayNames to be %v, got %v", customOptions.UseShortDayNames, customModel.useShortDayNames)
	}

	if customModel.dateRangeEnabled != (customOptions.EndMonth != nil) {
		t.Errorf("Expected dateRangeEnabled to be %v, got %v", customOptions.EndMonth != nil, customModel.dateRangeEnabled)
	}

	if customModel.justifyIndex != 1 { // 1 is the index for "center"
		t.Errorf("Expected justifyIndex to be %d, got %d", 1, customModel.justifyIndex)
	}
}

func TestValidateYear(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Valid year",
			input:   "2025",
			wantErr: false,
		},
		{
			name:    "Invalid year",
			input:   "not-a-year",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateYear(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateYear() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMonth(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Valid month",
			input:   "7",
			wantErr: false,
		},
		{
			name:    "Invalid month (not a number)",
			input:   "not-a-month",
			wantErr: true,
		},
		{
			name:    "Month out of range (too high)",
			input:   "13",
			wantErr: true,
		},
		{
			name:    "Month out of range (too low)",
			input:   "0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMonth(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMonth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateOptions(t *testing.T) {
	// Create a model with some initial values
	options := calendar.NewOptions()
	model := Initialize(&options)

	// Set some values in the model
	model.weekdayIndex = int(time.Sunday)
	model.showWeekNumbers = false
	model.showWeekends = false
	model.showComments = false
	model.useShortDayNames = true
	model.justifyIndex = 1 // "center"

	// Update the options
	err := model.updateOptions()
	if err != nil {
		t.Fatalf("updateOptions() error = %v", err)
	}

	// Check that the options were updated correctly
	if options.FirstDayOfWeek != time.Sunday {
		t.Errorf("Expected FirstDayOfWeek to be %v, got %v", time.Sunday, options.FirstDayOfWeek)
	}

	if options.ShowCalendarWeek != false {
		t.Errorf("Expected ShowCalendarWeek to be %v, got %v", false, options.ShowCalendarWeek)
	}

	if options.ShowWeekends != false {
		t.Errorf("Expected ShowWeekends to be %v, got %v", false, options.ShowWeekends)
	}

	if options.ShowComments != false {
		t.Errorf("Expected ShowComments to be %v, got %v", false, options.ShowComments)
	}

	if options.UseShortDayNames != true {
		t.Errorf("Expected UseShortDayNames to be %v, got %v", true, options.UseShortDayNames)
	}

	if options.Justify != "center" {
		t.Errorf("Expected Justify to be %v, got %v", "center", options.Justify)
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func ptrValue(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
