package layoutoptions

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

	docStyle = lipgloss.NewStyle().Margin(2, 5)
)

// Model represents the state of the layout options form
type Model struct {
	focusIndex       int
	weekdayOptions   []string
	weekdayIndex     int
	justifyOptions   []string
	justifyIndex     int
	showWeekNumbers  bool
	showWeekends     bool
	showComments     bool
	useShortDayNames bool
	done             bool
}

// NewModel creates a new layout options model
func NewModel() Model {
	return Model{
		focusIndex:       0,
		weekdayOptions:   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		weekdayIndex:     0, // Default to Monday
		justifyOptions:   []string{"left", "center", "right"},
		justifyIndex:     0, // Default to left
		showWeekNumbers:  true,
		showWeekends:     true,
		showComments:     true,
		useShortDayNames: false,
		done:             false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			// Cycle through options
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// Wrap around
			if m.focusIndex > 5 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 5
			}

		case "left", "right":
			// Handle selection for options
			switch m.focusIndex {
			case 0: // First day of week
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
			case 1: // Show week numbers
				m.showWeekNumbers = !m.showWeekNumbers
			case 2: // Show weekends
				m.showWeekends = !m.showWeekends
			case 3: // Show comments
				m.showComments = !m.showComments
			case 4: // Use short day names
				m.useShortDayNames = !m.useShortDayNames
			case 5: // Justification
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

		case "enter":
			m.done = true
			return m, nil
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	s := titleStyle.Render("Layout Options") + "\n\n"

	// First day of week
	weekdayLabel := optionLabelStyle.Render("First day of week:")
	weekdayOption := ""
	if m.focusIndex == 0 {
		weekdayOption = selectedItemStyle.Render("◀ " + m.weekdayOptions[m.weekdayIndex] + " ▶")
	} else {
		weekdayOption = unselectedItemStyle.Render(m.weekdayOptions[m.weekdayIndex])
	}
	s += weekdayLabel + " " + weekdayOption + "\n"

	// Show week numbers
	weekNumbersLabel := optionLabelStyle.Render("Show week numbers:")
	weekNumbersOption := ""
	if m.focusIndex == 1 {
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
	if m.focusIndex == 2 {
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
	if m.focusIndex == 3 {
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
	if m.focusIndex == 4 {
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
	if m.focusIndex == 5 {
		justifyOption = selectedItemStyle.Render("◀ " + m.justifyOptions[m.justifyIndex] + " ▶")
	} else {
		justifyOption = unselectedItemStyle.Render(m.justifyOptions[m.justifyIndex])
	}
	s += justifyLabel + " " + justifyOption + "\n\n"

	// Generate button
	s += focusedStyle.Render(" Generate Calendar ") + " (press Enter)\n"

	// Help text
	s += helpStyle.Render("\nUse tab/shift+tab or up/down to navigate, left/right to change selection, enter to generate")
	s += helpStyle.Render("\nPress Esc or Ctrl+C to quit")

	return docStyle.Render(s)
}

// Done returns whether the form is done
func (m Model) Done() bool {
	return m.done
}

// GetFirstDayOfWeek returns the selected first day of week
func (m Model) GetFirstDayOfWeek() time.Weekday {
	// Convert from UI index (Monday=0, Tuesday=1, etc.) to time.Weekday (Sunday=0, Monday=1, etc.)
	return time.Weekday((m.weekdayIndex + 1) % 7)
}

// GetShowWeekNumbers returns whether to show week numbers
func (m Model) GetShowWeekNumbers() bool {
	return m.showWeekNumbers
}

// GetShowWeekends returns whether to show weekends
func (m Model) GetShowWeekends() bool {
	return m.showWeekends
}

// GetShowComments returns whether to show comments
func (m Model) GetShowComments() bool {
	return m.showComments
}

// GetUseShortDayNames returns whether to use short day names
func (m Model) GetUseShortDayNames() bool {
	return m.useShortDayNames
}

// GetJustify returns the selected justification
func (m Model) GetJustify() string {
	return m.justifyOptions[m.justifyIndex]
}
