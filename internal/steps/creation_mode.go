package steps

import (
	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	tea "github.com/charmbracelet/bubbletea"
)

// CreationModeStep handles the creation mode selection
type CreationModeStep struct {
	list   ui.SelectList
	config *model.DBConfig
}

// NewCreationModeStep creates a new creation mode step
func NewCreationModeStep() *CreationModeStep {
	items := []ui.SelectItem{
		{
			Title:       "Typical Configuration",
			Description: "Create a database with minimal configuration using best practice defaults",
			Value:       string(model.CreationModeTypical),
		},
		{
			Title:       "Advanced Configuration",
			Description: "Create a database with full control over all configuration options",
			Value:       string(model.CreationModeAdvanced),
		},
	}

	return &CreationModeStep{
		list: ui.NewSelectList(items),
	}
}

// Init initializes the step
func (s *CreationModeStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.list.Reset()

	// Set cursor to current config value
	for i, item := range s.list.Items {
		if item.Value == string(config.CreationMode) {
			s.list.Cursor = i
			break
		}
	}

	return nil
}

// Update handles messages
func (s *CreationModeStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
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
func (s *CreationModeStep) View() string {
	return ui.SubtitleStyle.Render("Select the database creation mode:") + "\n\n" + s.list.View()
}

// Title returns the step title
func (s *CreationModeStep) Title() string {
	return "Database Creation Mode"
}

// Apply applies the step's changes to the config
func (s *CreationModeStep) Apply(config *model.DBConfig) {
	config.CreationMode = model.CreationMode(s.list.GetSelectedValue())
}

// ShouldSkip returns whether this step should be skipped
func (s *CreationModeStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
