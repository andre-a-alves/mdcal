package selection

import (
	"github.com/andre-a-alves/mdcal/cmd/steps"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	docStyle = lipgloss.NewStyle().Margin(2, 5)
)

// SelectionItem represents an item in the selection list
type SelectionItem struct {
	title, desc string
}

func (i SelectionItem) Title() string       { return i.title }
func (i SelectionItem) Description() string { return i.desc }
func (i SelectionItem) FilterValue() string { return i.title }

// Model represents the state of the selection list
type Model struct {
	list     list.Model
	selected string
	done     bool
	stepName string
}

// NewModel creates a new selection model
func NewModel(stepSchema steps.StepSchema) Model {
	items := []list.Item{}
	for _, option := range stepSchema.Options {
		items = append(items, SelectionItem{
			title: option.Title,
			desc:  option.Description,
		})
	}

	// Create a new list with default delegate and explicit size
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 30, 10)

	// Set the title from the step schema
	l.Title = stepSchema.Header

	// Set styles for better visibility
	l.SetShowTitle(false)        // Hide the title as we'll display it separately
	l.SetFilteringEnabled(false) // Disable filtering for simplicity

	return Model{
		list:     l,
		selected: "",
		done:     false,
		stepName: stepSchema.Name,
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
		case "enter":
			i, ok := m.list.SelectedItem().(SelectionItem)
			if ok {
				m.selected = i.Title()
				m.done = true
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		// Reserve space for the logo and question (approximately 10 lines)
		reservedHeight := 10

		// Get the frame size from docStyle
		h, v := docStyle.GetFrameSize()

		// Set the list size, reserving space for the logo and question
		m.list.SetSize(msg.Width-h, msg.Height-v-reservedHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	// Start with the logo
	output := titleStyle.Render("" +
		"                  __           __\n" +
		"   ____ ___  ____/ /________ _/ /\n" +
		"  / __ `__ \\/ __  / ___/ __ `/ / \n" +
		" / / / / / / /_/ / /__/ /_/ / /  \n" +
		"/_/ /_/ /_/\\__,_/\\___/\\__,_/_/   \n" +
		"📅 Markdown Calendar Generator")

	// Add the question
	output += "\n\n" + titleStyle.Render("What would you like to generate?") + "\n\n"

	// Get the currently selected index
	selectedIndex := m.list.Index()

	// Manually render each item in the list
	items := m.list.Items()
	for i, item := range items {
		selItem, ok := item.(SelectionItem)
		if !ok {
			continue
		}

		// Style based on whether this item is selected
		if i == selectedIndex {
			output += "▶ " + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2D7D9A")).Render(selItem.Title())
		} else {
			output += "  " + selItem.Title()
		}

		// Add the description
		if selItem.Description() != "" {
			output += " - " + lipgloss.NewStyle().Faint(true).Render(selItem.Description())
		}

		output += "\n"
	}

	// Add help text
	output += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓: Navigate • Enter: Select")

	return output
}

// Selected returns the selected item
func (m Model) Selected() string {
	return m.selected
}

// Done returns whether the selection is done
func (m Model) Done() bool {
	return m.done
}

// StepName returns the name of the step
func (m Model) StepName() string {
	return m.stepName
}
