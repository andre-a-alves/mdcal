package interactive

import (
	"fmt"
	"os"

	"github.com/andre-a-alves/mdcal/cmd/calendar"
	"github.com/andre-a-alves/mdcal/cmd/steps"
	"github.com/andre-a-alves/mdcal/cmd/ui/dateinput"
	"github.com/andre-a-alves/mdcal/cmd/ui/layoutoptions"
	"github.com/andre-a-alves/mdcal/cmd/ui/selection"
	tea "github.com/charmbracelet/bubbletea"
)

// Step represents the current step in the multi-step UI flow
type Step int

const (
	SelectionStep Step = iota
	DateInputStep
	LayoutOptionsStep
)

// MultiStepModel represents the state of the multi-step UI flow
type MultiStepModel struct {
	step           Step
	selectionModel selection.Model
	dateInputModel dateinput.Model
	layoutModel    layoutoptions.Model
	options        *calendar.Options
	quitting       bool
}

// InitializeMultiStep creates a new multi-step model
func InitializeMultiStep(options *calendar.Options) MultiStepModel {
	stepsData := steps.InitSteps()
	dateOptionsStep := stepsData.Steps["dateOptions"]

	return MultiStepModel{
		step:           SelectionStep,
		selectionModel: selection.NewModel(dateOptionsStep),
		options:        options,
	}
}

// Init initializes the model
func (m MultiStepModel) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model
func (m MultiStepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	// Handle updates based on the current step
	switch m.step {
	case SelectionStep:
		// Update the selection model
		newSelectionModel, selCmd := m.selectionModel.Update(msg)
		m.selectionModel = newSelectionModel.(selection.Model)
		cmd = selCmd

		// If selection is done, move to the date input step
		if m.selectionModel.Done() {
			selected := m.selectionModel.Selected()
			var inputMode dateinput.InputMode

			switch selected {
			case "Year":
				inputMode = dateinput.YearMode
			case "Month":
				inputMode = dateinput.MonthMode
			case "Range":
				inputMode = dateinput.RangeMode
			}

			m.dateInputModel = dateinput.NewModel(inputMode)
			m.step = DateInputStep
		}

	case DateInputStep:
		// Update the date input model
		newDateInputModel, dateCmd := m.dateInputModel.Update(msg)
		m.dateInputModel = newDateInputModel.(dateinput.Model)
		cmd = dateCmd

		// If date input is done, move to the layout options step
		if m.dateInputModel.Done() {
			m.layoutModel = layoutoptions.NewModel()
			m.step = LayoutOptionsStep
		}

	case LayoutOptionsStep:
		// Update the layout options model
		newLayoutModel, layoutCmd := m.layoutModel.Update(msg)
		m.layoutModel = newLayoutModel.(layoutoptions.Model)
		cmd = layoutCmd

		// If layout options are done, update the options and quit
		if m.layoutModel.Done() {
			m.updateOptions()
			return m, tea.Quit
		}
	}

	return m, cmd
}

// View renders the UI
func (m MultiStepModel) View() string {
	if m.quitting {
		return ""
	}

	// Render the appropriate view based on the current step
	switch m.step {
	case SelectionStep:
		return m.selectionModel.View()
	case DateInputStep:
		return m.dateInputModel.View()
	case LayoutOptionsStep:
		return m.layoutModel.View()
	}

	return "Unknown step"
}

// updateOptions updates the calendar options based on the user's selections
func (m *MultiStepModel) updateOptions() {
	// Update date options
	year, _ := m.dateInputModel.GetYear()
	m.options.Year = year

	month, _ := m.dateInputModel.GetMonth()
	m.options.Month = month

	endYear, _ := m.dateInputModel.GetEndYear()
	m.options.EndYear = endYear

	endMonth, _ := m.dateInputModel.GetEndMonth()
	m.options.EndMonth = endMonth

	// Update layout options
	m.options.FirstDayOfWeek = m.layoutModel.GetFirstDayOfWeek()
	m.options.ShowCalendarWeek = m.layoutModel.GetShowWeekNumbers()
	m.options.ShowWeekends = m.layoutModel.GetShowWeekends()
	m.options.ShowComments = m.layoutModel.GetShowComments()
	m.options.UseShortDayNames = m.layoutModel.GetUseShortDayNames()
	m.options.Justify = m.layoutModel.GetJustify()
}

// RunMultiStepMode runs the multi-step interactive mode
func RunMultiStepMode(options *calendar.Options) {
	p := tea.NewProgram(InitializeMultiStep(options), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running interactive mode: %v\n", err)
		os.Exit(1)
	}
}
