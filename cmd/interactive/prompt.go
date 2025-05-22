package interactive

import (
	"fmt"
	"github.com/andre-a-alves/mdcal/cmd/calendar"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#2D7D9A")).
			Padding(0, 1)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Padding(1, 0, 0, 2)

	optionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2D7D9A")).
				Width(30).
				PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#2D7D9A")).
				Padding(0, 1)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#DDDDDD")).
				Padding(0, 1)
)

// Model represents the state of the interactive form
type Model struct {
	inputs           []textinput.Model
	focusIndex       int
	cursorMode       textinput.CursorMode
	weekdayOptions   []string
	weekdayIndex     int
	justifyOptions   []string
	justifyIndex     int
	showWeekNumbers  bool
	showWeekends     bool
	showComments     bool
	useShortDayNames bool
	dateRangeEnabled bool
	options          *calendar.Options
	err              string
	success          string
	quitting         bool
}

// Initialize creates a new model with default values
func Initialize(options *calendar.Options) Model {
	m := Model{
		inputs:           make([]textinput.Model, 5),
		cursorMode:       textinput.CursorBlink,
		weekdayOptions:   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		weekdayIndex:     int(options.FirstDayOfWeek),
		justifyOptions:   []string{"left", "center", "right"},
		justifyIndex:     0, // Default to left
		showWeekNumbers:  options.ShowCalendarWeek,
		showWeekends:     options.ShowWeekends,
		showComments:     options.ShowComments,
		useShortDayNames: options.UseShortDayNames,
		dateRangeEnabled: options.EndMonth != nil,
		options:          options,
	}

	// Set justify index based on options
	for i, j := range m.justifyOptions {
		if j == options.Justify {
			m.justifyIndex = i
			break
		}
	}

	// Year input
	yearInput := textinput.New()
	yearInput.Placeholder = fmt.Sprintf("%d", options.Year)
	yearInput.Focus()
	yearInput.CharLimit = 4
	yearInput.Width = 6
	yearInput.Prompt = ""
	yearInput.PromptStyle = focusedStyle
	yearInput.TextStyle = focusedStyle
	yearInput.Validate = validateYear
	m.inputs[0] = yearInput

	// Month input
	monthInput := textinput.New()
	if options.Month != nil {
		monthInput.Placeholder = fmt.Sprintf("%d", *options.Month)
	} else {
		monthInput.Placeholder = "All"
	}
	monthInput.CharLimit = 2
	monthInput.Width = 6
	monthInput.Prompt = ""
	monthInput.PromptStyle = blurredStyle
	monthInput.Validate = validateMonth
	m.inputs[1] = monthInput

	// End Year input
	endYearInput := textinput.New()
	if options.EndYear != nil {
		endYearInput.Placeholder = fmt.Sprintf("%d", *options.EndYear)
	} else {
		endYearInput.Placeholder = fmt.Sprintf("%d", options.Year)
	}
	endYearInput.CharLimit = 4
	endYearInput.Width = 6
	endYearInput.Prompt = ""
	endYearInput.PromptStyle = blurredStyle
	endYearInput.Validate = validateYear
	m.inputs[2] = endYearInput

	// End Month input
	endMonthInput := textinput.New()
	if options.EndMonth != nil {
		endMonthInput.Placeholder = fmt.Sprintf("%d", *options.EndMonth)
	} else if options.Month != nil {
		endMonthInput.Placeholder = fmt.Sprintf("%d", *options.Month)
	} else {
		endMonthInput.Placeholder = "12"
	}
	endMonthInput.CharLimit = 2
	endMonthInput.Width = 6
	endMonthInput.Prompt = ""
	endMonthInput.PromptStyle = blurredStyle
	endMonthInput.Validate = validateMonth
	m.inputs[3] = endMonthInput

	return m
}

// validateYear checks if the input is a valid year
func validateYear(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.Atoi(s)
	return err
}

// validateMonth checks if the input is a valid month (1-12)
func validateMonth(s string) error {
	if s == "" {
		return nil
	}
	month, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}
	return nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			// Cycle through inputs
			oldFocusIndex := m.focusIndex

			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// Handle cycling between date options and layout options
			maxIndex := 9

			// Special handling for cycling between sections
			if !m.dateRangeEnabled {
				// If going up from Year (index 0), go to the bottom of layout options
				if (msg.String() == "up" || msg.String() == "shift+tab") && oldFocusIndex == 0 {
					m.focusIndex = maxIndex
				}

				// If going down from Month (index 1), go to layout options
				if (msg.String() == "down" || msg.String() == "tab") && oldFocusIndex == 1 {
					m.focusIndex = 4
				}

				// If going up from First day of week (index 4), go to Month
				if (msg.String() == "up" || msg.String() == "shift+tab") && oldFocusIndex == 4 {
					m.focusIndex = 1
				}

				// If going down from Cell justification (index 9), go to Year
				if (msg.String() == "down" || msg.String() == "tab") && oldFocusIndex == 9 {
					m.focusIndex = 0
				}

				// Skip date range inputs (indices 2 and 3)
				if m.focusIndex == 2 || m.focusIndex == 3 {
					if oldFocusIndex < 2 { // Coming from above
						m.focusIndex = 4 // Skip to layout options
					} else { // Coming from below
						m.focusIndex = 1 // Skip to month input
					}
				}
			}

			// Wrap around when reaching the boundaries
			if m.focusIndex > maxIndex {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = maxIndex
			}

			// Handle focus for text inputs
			for i := 0; i < len(m.inputs); i++ {
				// Skip end year and end month inputs if date range is not enabled
				if (i == 2 || i == 3) && !m.dateRangeEnabled {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = blurredStyle
					m.inputs[i].TextStyle = blurredStyle
					continue
				}

				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].Focus())
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = blurredStyle
					m.inputs[i].TextStyle = blurredStyle
				}
			}

		case "enter":
			if m.focusIndex == 9 {
				// Submit form
				err := m.updateOptions()
				if err != nil {
					m.err = err.Error()
				} else {
					m.success = "Calendar options updated successfully!"
					return m, tea.Quit
				}
			}

		case "left", "right":
			// Handle selection for non-text input options
			switch m.focusIndex {
			case 4: // First day of week
				if msg.String() == "left" {
					m.weekdayIndex--
					if m.weekdayIndex < 0 {
						m.weekdayIndex = len(m.weekdayOptions) - 1
					}
				} else {
					m.weekdayIndex++
					if m.weekdayIndex >= len(m.weekdayOptions) {
						m.weekdayIndex = 0
					}
				}
			case 5: // Show week numbers
				m.showWeekNumbers = !m.showWeekNumbers
			case 6: // Show weekends
				m.showWeekends = !m.showWeekends
			case 7: // Show comments
				m.showComments = !m.showComments
			case 8: // Use short day names
				m.useShortDayNames = !m.useShortDayNames
			case 9: // Justification
				if msg.String() == "left" {
					m.justifyIndex--
					if m.justifyIndex < 0 {
						m.justifyIndex = len(m.justifyOptions) - 1
					}
				} else {
					m.justifyIndex++
					if m.justifyIndex >= len(m.justifyOptions) {
						m.justifyIndex = 0
					}
				}
			}

		case "d":
			// Toggle date range
			if m.focusIndex == 1 || m.focusIndex == 2 || m.focusIndex == 3 {
				m.dateRangeEnabled = !m.dateRangeEnabled

				// Update focus based on date range state
				if m.dateRangeEnabled {
					// If date range is enabled and focus is on month, move to end year
					if m.focusIndex == 1 {
						m.focusIndex = 2
					}
				} else {
					// If date range is disabled and focus is on end year or end month, move to month
					if m.focusIndex == 2 || m.focusIndex == 3 {
						m.focusIndex = 1
					}
				}

				// Update input styles based on new focus
				for i := 0; i < len(m.inputs); i++ {
					// Skip end year and end month inputs if date range is not enabled
					if (i == 2 || i == 3) && !m.dateRangeEnabled {
						m.inputs[i].Blur()
						m.inputs[i].PromptStyle = blurredStyle
						m.inputs[i].TextStyle = blurredStyle
						continue
					}

					if i == m.focusIndex {
						cmds = append(cmds, m.inputs[i].Focus())
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
					} else {
						m.inputs[i].Blur()
						m.inputs[i].PromptStyle = blurredStyle
						m.inputs[i].TextStyle = blurredStyle
					}
				}
			}
		}
	}

	// Handle text input updates
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateInputs updates the text inputs based on user input
func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	for i := range m.inputs {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// updateOptions updates the calendar options based on the form values
func (m *Model) updateOptions() error {
	// Year
	if m.inputs[0].Value() != "" {
		year, err := strconv.Atoi(m.inputs[0].Value())
		if err != nil {
			return fmt.Errorf("invalid year: %v", err)
		}
		m.options.Year = year
	}

	// Month
	if m.inputs[1].Value() != "" {
		month, err := strconv.Atoi(m.inputs[1].Value())
		if err != nil {
			return fmt.Errorf("invalid month: %v", err)
		}
		if month < 1 || month > 12 {
			return fmt.Errorf("month must be between 1 and 12")
		}
		m.options.Month = &month
	} else {
		m.options.Month = nil
	}

	// Date range
	if m.dateRangeEnabled {
		// End Year
		endYear := m.options.Year
		if m.inputs[2].Value() != "" {
			var err error
			endYear, err = strconv.Atoi(m.inputs[2].Value())
			if err != nil {
				return fmt.Errorf("invalid end year: %v", err)
			}
		}
		m.options.EndYear = &endYear

		// End Month
		endMonth := 12
		if m.inputs[3].Value() != "" {
			var err error
			endMonth, err = strconv.Atoi(m.inputs[3].Value())
			if err != nil {
				return fmt.Errorf("invalid end month: %v", err)
			}
			if endMonth < 1 || endMonth > 12 {
				return fmt.Errorf("end month must be between 1 and 12")
			}
		}
		m.options.EndMonth = &endMonth
	} else {
		m.options.EndYear = nil
		m.options.EndMonth = nil
	}

	// First day of week
	m.options.FirstDayOfWeek = time.Weekday(m.weekdayIndex)

	// Boolean options
	m.options.ShowCalendarWeek = m.showWeekNumbers
	m.options.ShowWeekends = m.showWeekends
	m.options.ShowComments = m.showComments
	m.options.UseShortDayNames = m.useShortDayNames

	// Justification
	m.options.Justify = m.justifyOptions[m.justifyIndex]

	return nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	s := titleStyle.Render(""+
		"                  __           __\n"+
		"   ____ ___  ____/ /________ _/ /\n"+
		"  / __ `__ \\/ __  / ___/ __ `/ / \n"+
		" / / / / / / /_/ / /__/ /_/ / /  \n"+
		"/_/ /_/ /_/\\__,_/\\___/\\__,_/_/   \n"+
		"📅 Markdown Calendar Generator") + "\n\n"

	// Date section
	s += sectionTitleStyle.Render("Date Options") + "\n"

	// Year
	yearLabel := optionLabelStyle.Render("Year:")
	yearInput := m.inputs[0].View()
	if m.focusIndex == 0 {
		yearInput = focusedStyle.Render(yearInput)
	} else {
		yearInput = blurredStyle.Render(yearInput)
	}
	s += yearLabel + " " + yearInput + "\n"

	// Month
	monthLabel := optionLabelStyle.Render("Month (empty for whole year):")
	monthInput := m.inputs[1].View()
	if m.focusIndex == 1 {
		monthInput = focusedStyle.Render(monthInput)
	} else {
		monthInput = blurredStyle.Render(monthInput)
	}
	s += monthLabel + " " + monthInput + "\n"

	// Date Range
	dateRangeToggle := ""
	if m.dateRangeEnabled {
		dateRangeToggle = " [Press 'd' to disable date range]"
	} else {
		dateRangeToggle = " [Press 'd' to enable date range]"
	}

	if m.focusIndex == 1 || m.focusIndex == 2 || m.focusIndex == 3 {
		s += helpStyle.Render(dateRangeToggle) + "\n"
	}

	if m.dateRangeEnabled {
		// End Year
		endYearLabel := optionLabelStyle.Render("End Year:")
		endYearInput := m.inputs[2].View()
		if m.focusIndex == 2 {
			endYearInput = focusedStyle.Render(endYearInput)
		} else {
			endYearInput = blurredStyle.Render(endYearInput)
		}
		s += endYearLabel + " " + endYearInput + "\n"

		// End Month
		endMonthLabel := optionLabelStyle.Render("End Month:")
		endMonthInput := m.inputs[3].View()
		if m.focusIndex == 3 {
			endMonthInput = focusedStyle.Render(endMonthInput)
		} else {
			endMonthInput = blurredStyle.Render(endMonthInput)
		}
		s += endMonthLabel + " " + endMonthInput + "\n"
	}

	// Layout section
	s += sectionTitleStyle.Render("Layout Options") + "\n"

	// First day of week
	weekdayLabel := optionLabelStyle.Render("First day of week:")
	weekdayOption := ""
	if m.focusIndex == 4 {
		weekdayOption = selectedItemStyle.Render("◀ " + m.weekdayOptions[m.weekdayIndex] + " ▶")
	} else {
		weekdayOption = unselectedItemStyle.Render(m.weekdayOptions[m.weekdayIndex])
	}
	s += weekdayLabel + " " + weekdayOption + "\n"

	// Show week numbers
	weekNumbersLabel := optionLabelStyle.Render("Show week numbers:")
	weekNumbersOption := ""
	if m.focusIndex == 5 {
		if m.showWeekNumbers {
			weekNumbersOption = selectedItemStyle.Render("◀ Yes ▶")
		} else {
			weekNumbersOption = selectedItemStyle.Render("◀ No ▶")
		}
	} else {
		if m.showWeekNumbers {
			weekNumbersOption = unselectedItemStyle.Render("Yes")
		} else {
			weekNumbersOption = unselectedItemStyle.Render("No")
		}
	}
	s += weekNumbersLabel + " " + weekNumbersOption + "\n"

	// Show weekends
	weekendsLabel := optionLabelStyle.Render("Show weekends:")
	weekendsOption := ""
	if m.focusIndex == 6 {
		if m.showWeekends {
			weekendsOption = selectedItemStyle.Render("◀ Yes ▶")
		} else {
			weekendsOption = selectedItemStyle.Render("◀ No ▶")
		}
	} else {
		if m.showWeekends {
			weekendsOption = unselectedItemStyle.Render("Yes")
		} else {
			weekendsOption = unselectedItemStyle.Render("No")
		}
	}
	s += weekendsLabel + " " + weekendsOption + "\n"

	// Show comments
	commentsLabel := optionLabelStyle.Render("Show comments column:")
	commentsOption := ""
	if m.focusIndex == 7 {
		if m.showComments {
			commentsOption = selectedItemStyle.Render("◀ Yes ▶")
		} else {
			commentsOption = selectedItemStyle.Render("◀ No ▶")
		}
	} else {
		if m.showComments {
			commentsOption = unselectedItemStyle.Render("Yes")
		} else {
			commentsOption = unselectedItemStyle.Render("No")
		}
	}
	s += commentsLabel + " " + commentsOption + "\n"

	// Use short day names
	shortDayNamesLabel := optionLabelStyle.Render("Use short day names:")
	shortDayNamesOption := ""
	if m.focusIndex == 8 {
		if m.useShortDayNames {
			shortDayNamesOption = selectedItemStyle.Render("◀ Yes ▶")
		} else {
			shortDayNamesOption = selectedItemStyle.Render("◀ No ▶")
		}
	} else {
		if m.useShortDayNames {
			shortDayNamesOption = unselectedItemStyle.Render("Yes")
		} else {
			shortDayNamesOption = unselectedItemStyle.Render("No")
		}
	}
	s += shortDayNamesLabel + " " + shortDayNamesOption + "\n"

	// Justification
	justifyLabel := optionLabelStyle.Render("Cell justification:")
	justifyOption := ""
	if m.focusIndex == 9 {
		justifyOption = selectedItemStyle.Render("◀ " + m.justifyOptions[m.justifyIndex] + " ▶")
	} else {
		justifyOption = unselectedItemStyle.Render(m.justifyOptions[m.justifyIndex])
	}
	s += justifyLabel + " " + justifyOption + "\n\n"

	// Submit button
	if m.focusIndex == 9 {
		s += focusedStyle.Render(" Generate Calendar ") + " (press Enter)\n"
	} else {
		s += blurredStyle.Render(" Generate Calendar ") + "\n"
	}

	// Help text
	s += helpStyle.Render("\nUse tab/shift+tab or up/down to navigate, left/right to change selection, enter to submit\n")
	s += helpStyle.Render("Press Esc or Ctrl+C to quit\n")

	// Error message
	if m.err != "" {
		s += "\n" + errorStyle.Render(m.err)
	}

	// Success message
	if m.success != "" {
		s += "\n" + successStyle.Render(m.success)
	}

	return s
}

// RunInteractiveMode runs the interactive mode using bubbletea
func RunInteractiveMode(options *calendar.Options) {
	// Use the new multi-step UI flow
	RunMultiStepMode(options)
}
