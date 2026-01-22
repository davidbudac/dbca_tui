package steps

import (
	"strings"

	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DataVaultStep handles Oracle Data Vault configuration
type DataVaultStep struct {
	config          *model.DBConfig
	inputs          []textinput.Model
	focusIndex      int
	enableDataVault bool
	err             string
}

const (
	dvIdxOwner = iota
	dvIdxAccountManager
)

// NewDataVaultStep creates a new Data Vault step
func NewDataVaultStep() *DataVaultStep {
	s := &DataVaultStep{
		inputs: make([]textinput.Model, 2),
	}

	// Data Vault Owner
	s.inputs[dvIdxOwner] = textinput.New()
	s.inputs[dvIdxOwner].Placeholder = "C##DVOWNER"
	s.inputs[dvIdxOwner].CharLimit = 30

	// Data Vault Account Manager
	s.inputs[dvIdxAccountManager] = textinput.New()
	s.inputs[dvIdxAccountManager].Placeholder = "C##DVACCTMGR"
	s.inputs[dvIdxAccountManager].CharLimit = 30

	return s
}

// Init initializes the step
func (s *DataVaultStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""
	s.enableDataVault = config.EnableDataVault

	// Set values from config
	s.inputs[dvIdxOwner].SetValue(config.DataVaultOwner)
	s.inputs[dvIdxAccountManager].SetValue(config.DataVaultAccountManager)

	// Reset focus
	for i := range s.inputs {
		s.inputs[i].Blur()
	}

	return nil
}

// Update handles messages
func (s *DataVaultStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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

		case "d", "D":
			// Toggle Data Vault
			if s.focusIndex == 0 || !s.enableDataVault {
				s.enableDataVault = !s.enableDataVault
				if s.enableDataVault {
					s.focusIndex = 1
					s.inputs[dvIdxOwner].Focus()
				} else {
					for i := range s.inputs {
						s.inputs[i].Blur()
					}
					s.focusIndex = 0
				}
				return s, wizard.StepStay, textinput.Blink
			}
		}
	}

	// Update the focused text input
	if s.enableDataVault && s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		inputIdx := s.focusIndex - 1
		var cmd tea.Cmd
		s.inputs[inputIdx], cmd = s.inputs[inputIdx].Update(msg)
		return s, wizard.StepStay, cmd
	}

	return s, wizard.StepStay, nil
}

func (s *DataVaultStep) nextField() {
	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Blur()
	}

	maxFields := 0 // Just the toggle
	if s.enableDataVault {
		maxFields = len(s.inputs) // toggle + inputs
	}

	s.focusIndex++
	if s.focusIndex > maxFields {
		s.focusIndex = 0
	}

	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Focus()
	}
}

func (s *DataVaultStep) prevField() {
	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Blur()
	}

	maxFields := 0
	if s.enableDataVault {
		maxFields = len(s.inputs)
	}

	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxFields
	}

	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Focus()
	}
}

func (s *DataVaultStep) validate() bool {
	s.err = ""

	if s.enableDataVault {
		owner := strings.TrimSpace(s.inputs[dvIdxOwner].Value())
		if owner == "" {
			s.err = "Data Vault Owner is required"
			return false
		}

		acctMgr := strings.TrimSpace(s.inputs[dvIdxAccountManager].Value())
		if acctMgr == "" {
			s.err = "Data Vault Account Manager is required"
			return false
		}
	}

	return true
}

// View renders the step
func (s *DataVaultStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure Oracle Data Vault:") + "\n\n")

	b.WriteString(ui.SubtitleStyle.Render("Oracle Data Vault provides controls to prevent unauthorized access") + "\n")
	b.WriteString(ui.SubtitleStyle.Render("to data by privileged database users.") + "\n\n")

	// Enable Data Vault toggle
	checkbox := ui.UncheckedStyle.String()
	if s.enableDataVault {
		checkbox = ui.CheckedStyle.String()
	}
	dvStyle := ui.NormalItemStyle
	if s.focusIndex == 0 {
		dvStyle = ui.SelectedItemStyle
	}
	b.WriteString(checkbox + " " + dvStyle.Render("Enable Oracle Data Vault") + "\n")
	b.WriteString(ui.SubtitleStyle.Render("    Press 'd' to toggle") + "\n\n")

	if s.enableDataVault {
		// Data Vault Owner
		b.WriteString(s.renderField("Data Vault Owner", s.inputs[dvIdxOwner], 1) + "\n")

		// Data Vault Account Manager
		b.WriteString(s.renderField("Data Vault Account Manager", s.inputs[dvIdxAccountManager], 2) + "\n")
	}

	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *DataVaultStep) renderField(label string, input textinput.Model, fieldIndex int) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if s.focusIndex == fieldIndex {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	)
}

// Title returns the step title
func (s *DataVaultStep) Title() string {
	return "Data Vault Configuration"
}

// Apply applies the step's changes to the config
func (s *DataVaultStep) Apply(config *model.DBConfig) {
	config.EnableDataVault = s.enableDataVault

	if s.enableDataVault {
		config.DataVaultOwner = strings.TrimSpace(s.inputs[dvIdxOwner].Value())
		config.DataVaultAccountManager = strings.TrimSpace(s.inputs[dvIdxAccountManager].Value())
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *DataVaultStep) ShouldSkip(config *model.DBConfig) bool {
	// Only show in advanced mode for create operation
	return config.Operation != model.OperationCreate || config.CreationMode != model.CreationModeAdvanced
}
