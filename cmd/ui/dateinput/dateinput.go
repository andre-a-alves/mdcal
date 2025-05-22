package dateinput

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
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

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				Padding(1, 0, 0, 2)

	docStyle = lipgloss.NewStyle().Margin(2, 5)
)

// InputMode represents the type of date input
type InputMode string

const (
	YearMode  InputMode = "Year"
	MonthMode InputMode = "Month"
	RangeMode InputMode = "Range"
)

// Model represents the state of the date input form
type Model struct {
	inputs     []textinput.Model
	focusIndex int
	cursorMode textinput.CursorMode
	mode       InputMode
	err        string
	done       bool
}

// NewModel creates a new date input model
func NewModel(mode InputMode) Model {
	m := Model{
		inputs:     make([]textinput.Model, 4), // year, month, endYear, endMonth
		cursorMode: textinput.CursorBlink,
		mode:       mode,
		done:       false,
	}

	// Year input
	yearInput := textinput.New()
	yearInput.Placeholder = fmt.Sprintf("%d", time.Now().Year())
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
	monthInput.Placeholder = fmt.Sprintf("%d", time.Now().Month())
	monthInput.CharLimit = 2
	monthInput.Width = 6
	monthInput.Prompt = ""
	monthInput.PromptStyle = blurredStyle
	monthInput.Validate = validateMonth
	m.inputs[1] = monthInput

	// End Year input
	endYearInput := textinput.New()
	endYearInput.Placeholder = fmt.Sprintf("%d", time.Now().Year())
	endYearInput.CharLimit = 4
	endYearInput.Width = 6
	endYearInput.Prompt = ""
	endYearInput.PromptStyle = blurredStyle
	endYearInput.Validate = validateYear
	m.inputs[2] = endYearInput

	// End Month input
	endMonthInput := textinput.New()
	endMonthInput.Placeholder = fmt.Sprintf("%d", time.Now().Month())
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
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			// Cycle through inputs
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// Handle cycling based on mode
			var maxIndex int
			switch m.mode {
			case YearMode:
				maxIndex = 0 // Only year input
			case MonthMode:
				maxIndex = 1 // Year and month inputs
			case RangeMode:
				maxIndex = 3 // All inputs
			}

			// Wrap around
			if m.focusIndex > maxIndex {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = maxIndex
			}

			// Update focus
			for i := 0; i <= maxIndex; i++ {
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
			// All fields can be empty, defaults will be used
			valid := true

			// For Range mode, if end year is empty, use start year
			if m.mode == RangeMode && m.inputs[2].Value() == "" && m.inputs[0].Value() != "" {
				m.inputs[2].SetValue(m.inputs[0].Value())
			}

			if valid {
				m.done = true
				m.err = ""
				return m, nil
			}
		}
	}

	// Handle text input updates
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	// If in Range mode and start year changes, update end year placeholder
	if m.mode == RangeMode && m.focusIndex == 0 && m.inputs[0].Value() != "" {
		m.inputs[2].Placeholder = m.inputs[0].Value()
	}

	return m, tea.Batch(cmds...)
}

// updateInputs updates the text inputs based on user input
func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	// Only update inputs that are relevant to the current mode
	var maxIndex int
	switch m.mode {
	case YearMode:
		maxIndex = 0
	case MonthMode:
		maxIndex = 1
	case RangeMode:
		maxIndex = 3
	}

	for i := 0; i <= maxIndex; i++ {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	s := titleStyle.Render("Date Options") + "\n\n"

	// Year input (always shown)
	s += "Year: " + m.inputs[0].View() + "\n"

	// Month input (shown for Month and Range modes)
	if m.mode == MonthMode || m.mode == RangeMode {
		s += "Month: " + m.inputs[1].View() + "\n"
	}

	// End Year and End Month inputs (shown for Range mode)
	if m.mode == RangeMode {
		s += "End Year: " + m.inputs[2].View() + "\n"
		s += "End Month: " + m.inputs[3].View() + "\n"
	}

	// Help text
	s += "\n" + helpStyle.Render("Use tab/shift+tab or up/down to navigate, enter to continue")
	s += "\n" + helpStyle.Render("Press Esc or Ctrl+C to quit")

	// Error message
	if m.err != "" {
		s += "\n\n" + errorStyle.Render(m.err)
	}

	return docStyle.Render(s)
}

// Done returns whether the form is done
func (m Model) Done() bool {
	return m.done
}

// GetYear returns the entered year
func (m Model) GetYear() (int, error) {
	if m.inputs[0].Value() == "" {
		return time.Now().Year(), nil
	}
	return strconv.Atoi(m.inputs[0].Value())
}

// GetMonth returns the entered month
func (m Model) GetMonth() (*int, error) {
	if m.inputs[1].Value() == "" {
		if m.mode == YearMode {
			return nil, nil // Whole year
		}
		month := int(time.Now().Month())
		return &month, nil
	}
	month, err := strconv.Atoi(m.inputs[1].Value())
	if err != nil {
		return nil, err
	}
	return &month, nil
}

// GetEndYear returns the entered end year
func (m Model) GetEndYear() (*int, error) {
	if m.mode != RangeMode {
		return nil, nil
	}

	if m.inputs[2].Value() == "" {
		// Use start year if end year is empty
		year, err := m.GetYear()
		if err != nil {
			return nil, err
		}
		return &year, nil
	}

	endYear, err := strconv.Atoi(m.inputs[2].Value())
	if err != nil {
		return nil, err
	}
	return &endYear, nil
}

// GetEndMonth returns the entered end month
func (m Model) GetEndMonth() (*int, error) {
	if m.mode != RangeMode {
		return nil, nil
	}

	if m.inputs[3].Value() == "" {
		// Use current month if end month is empty
		month := int(time.Now().Month())
		return &month, nil
	}

	endMonth, err := strconv.Atoi(m.inputs[3].Value())
	if err != nil {
		return nil, err
	}
	return &endMonth, nil
}
