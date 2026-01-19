package steps

import (
	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	tea "github.com/charmbracelet/bubbletea"
)

// TemplateStep handles the database template selection
type TemplateStep struct {
	list   ui.SelectList
	config *model.DBConfig
}

// NewTemplateStep creates a new template step
func NewTemplateStep() *TemplateStep {
	items := []ui.SelectItem{
		{
			Title:       "General Purpose / Transaction Processing",
			Description: "A pre-configured database template optimized for general purpose or OLTP workloads",
			Value:       string(model.TemplateGeneralPurpose),
		},
		{
			Title:       "Data Warehouse",
			Description: "A pre-configured database template optimized for data warehousing workloads",
			Value:       string(model.TemplateDataWarehouse),
		},
		{
			Title:       "Custom Database",
			Description: "Create a database with custom configuration (no template)",
			Value:       string(model.TemplateCustom),
		},
	}

	return &TemplateStep{
		list: ui.NewSelectList(items),
	}
}

// Init initializes the step
func (s *TemplateStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.list.Reset()

	for i, item := range s.list.Items {
		if item.Value == string(config.TemplateName) {
			s.list.Cursor = i
			break
		}
	}

	return nil
}

// Update handles messages
func (s *TemplateStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, wizard.StepBack, nil
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
func (s *TemplateStep) View() string {
	return ui.SubtitleStyle.Render("Select a database template:") + "\n\n" + s.list.View()
}

// Title returns the step title
func (s *TemplateStep) Title() string {
	return "Database Template"
}

// Apply applies the step's changes to the config
func (s *TemplateStep) Apply(config *model.DBConfig) {
	config.TemplateName = model.DatabaseTemplate(s.list.GetSelectedValue())

	// Set database type based on template
	switch config.TemplateName {
	case model.TemplateGeneralPurpose:
		config.DatabaseType = model.DatabaseTypeMultipurpose
	case model.TemplateDataWarehouse:
		config.DatabaseType = model.DatabaseTypeDataWarehouse
	default:
		config.DatabaseType = model.DatabaseTypeMultipurpose
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *TemplateStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
