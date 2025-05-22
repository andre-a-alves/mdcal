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

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = stepSchema.Header

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
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	return docStyle.Render(m.list.View())
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
