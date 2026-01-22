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

// NetworkStep handles network/listener configuration
type NetworkStep struct {
	config          *model.DBConfig
	inputs          []textinput.Model
	focusIndex      int
	createListener  bool
	err             string
}

const (
	netIdxListenerName = iota
	netIdxListenerPort
)

// NewNetworkStep creates a new network step
func NewNetworkStep() *NetworkStep {
	s := &NetworkStep{
		inputs: make([]textinput.Model, 2),
	}

	// Listener Name
	s.inputs[netIdxListenerName] = textinput.New()
	s.inputs[netIdxListenerName].Placeholder = "LISTENER"
	s.inputs[netIdxListenerName].CharLimit = 30

	// Listener Port
	s.inputs[netIdxListenerPort] = textinput.New()
	s.inputs[netIdxListenerPort].Placeholder = "1521"
	s.inputs[netIdxListenerPort].CharLimit = 5

	return s
}

// Init initializes the step
func (s *NetworkStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""
	s.createListener = config.CreateNewListener

	// Set values from config
	s.inputs[netIdxListenerName].SetValue(config.ListenerName)
	s.inputs[netIdxListenerPort].SetValue(strconv.Itoa(config.ListenerPort))

	// Reset focus
	for i := range s.inputs {
		s.inputs[i].Blur()
	}
	s.inputs[0].Focus()

	return textinput.Blink
}

// Update handles messages
func (s *NetworkStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, wizard.StepBack, nil

		case "tab", "down":
			s.nextField()
			return s, wizard.StepStay, nil

		case "shift+tab", "up":
			s.prevField()
			return s, wizard.StepStay, nil

		case "enter":
			if s.validate() {
				return s, wizard.StepContinue, nil
			}
			return s, wizard.StepStay, nil

		case "c", "C":
			// Toggle create new listener
			if s.focusIndex == len(s.inputs) {
				s.createListener = !s.createListener
			}
		}
	}

	// Update the focused text input
	if s.focusIndex < len(s.inputs) {
		var cmd tea.Cmd
		s.inputs[s.focusIndex], cmd = s.inputs[s.focusIndex].Update(msg)
		return s, wizard.StepStay, cmd
	}

	return s, wizard.StepStay, nil
}

func (s *NetworkStep) nextField() {
	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Blur()
	}

	s.focusIndex++
	if s.focusIndex > len(s.inputs) {
		s.focusIndex = 0
	}

	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *NetworkStep) prevField() {
	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Blur()
	}

	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = len(s.inputs)
	}

	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *NetworkStep) validate() bool {
	s.err = ""

	listenerName := strings.TrimSpace(s.inputs[netIdxListenerName].Value())
	if listenerName == "" {
		s.err = "Listener name is required"
		return false
	}

	portStr := strings.TrimSpace(s.inputs[netIdxListenerPort].Value())
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		s.err = "Port must be between 1 and 65535"
		return false
	}

	return true
}

// View renders the step
func (s *NetworkStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure network listener:") + "\n\n")

	// Listener Name
	b.WriteString(s.renderField("Listener Name", s.inputs[netIdxListenerName], 0) + "\n")

	// Listener Port
	b.WriteString(s.renderField("Listener Port", s.inputs[netIdxListenerPort], 1) + "\n")

	// Create new listener toggle
	checkbox := ui.UncheckedStyle.String()
	if s.createListener {
		checkbox = ui.CheckedStyle.String()
	}
	createStyle := ui.NormalItemStyle
	if s.focusIndex == len(s.inputs) {
		createStyle = ui.SelectedItemStyle
	}
	b.WriteString("\n" + checkbox + " " + createStyle.Render("Create new listener (if not exists)") + "\n")
	b.WriteString(ui.SubtitleStyle.Render("    Press 'c' to toggle") + "\n")

	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *NetworkStep) renderField(label string, input textinput.Model, index int) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if s.focusIndex == index {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	) + "\n"
}

// Title returns the step title
func (s *NetworkStep) Title() string {
	return "Network Configuration"
}

// Apply applies the step's changes to the config
func (s *NetworkStep) Apply(config *model.DBConfig) {
	config.ListenerName = strings.TrimSpace(s.inputs[netIdxListenerName].Value())
	port, _ := strconv.Atoi(strings.TrimSpace(s.inputs[netIdxListenerPort].Value()))
	config.ListenerPort = port
	config.CreateNewListener = s.createListener
}

// ShouldSkip returns whether this step should be skipped
func (s *NetworkStep) ShouldSkip(config *model.DBConfig) bool {
	// Skip for delete operation or in typical mode
	return config.Operation != model.OperationCreate || config.CreationMode == model.CreationModeTypical
}
