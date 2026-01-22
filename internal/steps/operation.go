package steps

import (
	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	tea "github.com/charmbracelet/bubbletea"
)

// OperationStep handles the operation selection (Create/Delete)
type OperationStep struct {
	list   ui.SelectList
	config *model.DBConfig
}

// NewOperationStep creates a new operation step
func NewOperationStep() *OperationStep {
	items := []ui.SelectItem{
		{
			Title:       "Create a Database",
			Description: "Create a new Oracle database with the DBCA wizard",
			Value:       string(model.OperationCreate),
		},
		{
			Title:       "Delete a Database",
			Description: "Generate command to delete an existing Oracle database",
			Value:       string(model.OperationDelete),
		},
	}

	return &OperationStep{
		list: ui.NewSelectList(items),
	}
}

// Init initializes the step
func (s *OperationStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.list.Reset()

	// Set cursor to current config value
	for i, item := range s.list.Items {
		if item.Value == string(config.Operation) {
			s.list.Cursor = i
			break
		}
	}

	return nil
}

// Update handles messages
func (s *OperationStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return s, wizard.StepQuit, nil
		case "enter", " ":
			s.list.Update(msg)
			if s.list.IsSelected() {
				return s, wizard.StepContinue, nil
			}
		default:
			s.list.Update(msg)
		}
	}

	return s, wizard.StepStay, nil
}

// View renders the step
func (s *OperationStep) View() string {
	return ui.SubtitleStyle.Render("What would you like to do?") + "\n\n" + s.list.View()
}

// Title returns the step title
func (s *OperationStep) Title() string {
	return "Select Operation"
}

// Apply applies the step's changes to the config
func (s *OperationStep) Apply(config *model.DBConfig) {
	config.Operation = model.Operation(s.list.GetSelectedValue())
}

// ShouldSkip returns whether this step should be skipped
func (s *OperationStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
