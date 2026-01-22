package steps

import (
	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	tea "github.com/charmbracelet/bubbletea"
)

// DeploymentStep handles the deployment type selection
type DeploymentStep struct {
	list   ui.SelectList
	config *model.DBConfig
}

// NewDeploymentStep creates a new deployment step
func NewDeploymentStep() *DeploymentStep {
	items := []ui.SelectItem{
		{
			Title:       "Oracle Single Instance Database",
			Description: "A single database instance running on one server",
			Value:       string(model.DeploymentSingleInstance),
		},
		{
			Title:       "Oracle RAC Database",
			Description: "A clustered database with multiple instances across multiple nodes",
			Value:       string(model.DeploymentRAC),
		},
		{
			Title:       "Oracle RAC One Node Database",
			Description: "A single instance on one node with failover capability to other cluster nodes",
			Value:       string(model.DeploymentRACOneNode),
		},
	}

	return &DeploymentStep{
		list: ui.NewSelectList(items),
	}
}

// Init initializes the step
func (s *DeploymentStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.list.Reset()

	for i, item := range s.list.Items {
		if item.Value == string(config.DeploymentType) {
			s.list.Cursor = i
			break
		}
	}

	return nil
}

// Update handles messages
func (s *DeploymentStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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
func (s *DeploymentStep) View() string {
	return ui.SubtitleStyle.Render("Select the database deployment type:") + "\n\n" + s.list.View()
}

// Title returns the step title
func (s *DeploymentStep) Title() string {
	return "Deployment Type"
}

// Apply applies the step's changes to the config
func (s *DeploymentStep) Apply(config *model.DBConfig) {
	config.DeploymentType = model.DeploymentType(s.list.GetSelectedValue())
}

// ShouldSkip returns whether this step should be skipped
func (s *DeploymentStep) ShouldSkip(config *model.DBConfig) bool {
	return config.Operation != model.OperationCreate
}
