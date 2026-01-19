package steps

import (
	"strconv"
	"strings"

	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ManagementStep handles Enterprise Manager configuration
type ManagementStep struct {
	config     *model.DBConfig
	emList     ui.SelectList
	portInput  textinput.Model
	agentInput textinput.Model
	phase      int // 0=EM type, 1=port/agent config
	focusIndex int
	err        string
}

// NewManagementStep creates a new management step
func NewManagementStep() *ManagementStep {
	emItems := []ui.SelectItem{
		{
			Title:       "Do not configure Enterprise Manager",
			Description: "Skip Enterprise Manager configuration",
			Value:       string(model.EMConfigNone),
		},
		{
			Title:       "Configure Enterprise Manager Database Express",
			Description: "Built-in web-based database management (port 5500)",
			Value:       string(model.EMConfigDBExpress),
		},
		{
			Title:       "Register with Enterprise Manager Cloud Control",
			Description: "Register with an existing Cloud Control installation",
			Value:       string(model.EMConfigCentral),
		},
	}

	s := &ManagementStep{
		emList: ui.NewSelectList(emItems),
	}

	s.portInput = textinput.New()
	s.portInput.Placeholder = "5500"
	s.portInput.CharLimit = 5

	s.agentInput = textinput.New()
	s.agentInput.Placeholder = "hostname:3872"
	s.agentInput.CharLimit = 256

	return s
}

// Init initializes the step
func (s *ManagementStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.phase = 0
	s.focusIndex = 0
	s.err = ""
	s.emList.Reset()

	// Set cursor to current config value
	for i, item := range s.emList.Items {
		if item.Value == string(config.EMConfiguration) {
			s.emList.Cursor = i
			break
		}
	}

	s.portInput.SetValue(strconv.Itoa(config.EMPort))
	s.agentInput.SetValue(config.CloudControlAgent)

	return nil
}

// Update handles messages
func (s *ManagementStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.phase > 0 {
				s.phase = 0
				s.portInput.Blur()
				s.agentInput.Blur()
				return s, wizard.StepStay, nil
			}
			return s, wizard.StepBack, nil
		}
	}

	if s.phase == 0 {
		return s.updateEMSelection(msg)
	}
	return s.updatePortConfig(msg)
}

func (s *ManagementStep) updateEMSelection(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			s.emList.Update(msg)
			if s.emList.IsSelected() {
				emConfig := model.EMConfiguration(s.emList.GetSelectedValue())
				if emConfig == model.EMConfigNone {
					return s, wizard.StepContinue, nil
				}
				s.phase = 1
				s.focusIndex = 0
				if emConfig == model.EMConfigDBExpress {
					s.portInput.Focus()
				} else {
					s.agentInput.Focus()
				}
				return s, wizard.StepStay, textinput.Blink
			}
		default:
			s.emList.Update(msg)
		}
	}
	return s, wizard.StepStay, nil
}

func (s *ManagementStep) updatePortConfig(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	emConfig := model.EMConfiguration(s.emList.GetSelectedValue())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "shift+tab", "up":
			// For Central config, allow switching between fields
			if emConfig == model.EMConfigCentral {
				if s.focusIndex == 0 {
					s.agentInput.Blur()
					s.focusIndex = 1
					s.portInput.Focus()
				} else {
					s.portInput.Blur()
					s.focusIndex = 0
					s.agentInput.Focus()
				}
			}
			return s, wizard.StepStay, nil

		case "enter":
			if s.validate(emConfig) {
				return s, wizard.StepContinue, nil
			}
			return s, wizard.StepStay, nil
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	if emConfig == model.EMConfigDBExpress {
		s.portInput, cmd = s.portInput.Update(msg)
	} else {
		if s.focusIndex == 0 {
			s.agentInput, cmd = s.agentInput.Update(msg)
		} else {
			s.portInput, cmd = s.portInput.Update(msg)
		}
	}
	return s, wizard.StepStay, cmd
}

func (s *ManagementStep) validate(emConfig model.EMConfiguration) bool {
	s.err = ""

	if emConfig == model.EMConfigDBExpress || emConfig == model.EMConfigCentral {
		portStr := strings.TrimSpace(s.portInput.Value())
		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			s.err = "Port must be between 1 and 65535"
			return false
		}
	}

	if emConfig == model.EMConfigCentral {
		agent := strings.TrimSpace(s.agentInput.Value())
		if agent == "" {
			s.err = "Cloud Control agent URL is required"
			return false
		}
	}

	return true
}

// View renders the step
func (s *ManagementStep) View() string {
	var b strings.Builder

	if s.phase == 0 {
		b.WriteString(ui.SubtitleStyle.Render("Configure database management options:") + "\n\n")
		b.WriteString(s.emList.View())
	} else {
		emConfig := model.EMConfiguration(s.emList.GetSelectedValue())
		emItem := s.emList.GetSelectedItem()

		b.WriteString(ui.SubtitleStyle.Render(emItem.Title) + "\n\n")

		if emConfig == model.EMConfigDBExpress {
			b.WriteString(s.renderField("HTTPS Port", s.portInput, true) + "\n")
			b.WriteString(ui.SubtitleStyle.Render("    Access URL: https://hostname:PORT/em") + "\n")
		} else {
			b.WriteString(s.renderField("Cloud Control Agent URL", s.agentInput, s.focusIndex == 0) + "\n")
			b.WriteString(s.renderField("Agent Port", s.portInput, s.focusIndex == 1) + "\n")
		}

		if s.err != "" {
			b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
		}

		b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))
	}

	return b.String()
}

func (s *ManagementStep) renderField(label string, input textinput.Model, focused bool) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if focused {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	) + "\n"
}

// Title returns the step title
func (s *ManagementStep) Title() string {
	return "Management Options"
}

// Apply applies the step's changes to the config
func (s *ManagementStep) Apply(config *model.DBConfig) {
	config.EMConfiguration = model.EMConfiguration(s.emList.GetSelectedValue())

	if config.EMConfiguration != model.EMConfigNone {
		port, _ := strconv.Atoi(strings.TrimSpace(s.portInput.Value()))
		config.EMPort = port
	}

	if config.EMConfiguration == model.EMConfigCentral {
		config.CloudControlAgent = strings.TrimSpace(s.agentInput.Value())
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *ManagementStep) ShouldSkip(config *model.DBConfig) bool {
	// Skip in typical mode
	return config.CreationMode == model.CreationModeTypical
}
