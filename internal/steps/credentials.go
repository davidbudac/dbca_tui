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

// CredentialsStep handles database credentials configuration
type CredentialsStep struct {
	config            *model.DBConfig
	inputs            []textinput.Model
	focusIndex        int
	useCommonPassword bool
	err               string
}

const (
	credIdxCommon = iota
	credIdxSys
	credIdxSystem
	credIdxPDBAdmin
)

// NewCredentialsStep creates a new credentials step
func NewCredentialsStep() *CredentialsStep {
	s := &CredentialsStep{
		inputs: make([]textinput.Model, 4),
	}

	// Common Password
	s.inputs[credIdxCommon] = textinput.New()
	s.inputs[credIdxCommon].Placeholder = "Enter password for all accounts"
	s.inputs[credIdxCommon].EchoMode = textinput.EchoPassword
	s.inputs[credIdxCommon].EchoCharacter = '*'
	s.inputs[credIdxCommon].CharLimit = 30

	// SYS Password
	s.inputs[credIdxSys] = textinput.New()
	s.inputs[credIdxSys].Placeholder = "SYS password"
	s.inputs[credIdxSys].EchoMode = textinput.EchoPassword
	s.inputs[credIdxSys].EchoCharacter = '*'
	s.inputs[credIdxSys].CharLimit = 30

	// SYSTEM Password
	s.inputs[credIdxSystem] = textinput.New()
	s.inputs[credIdxSystem].Placeholder = "SYSTEM password"
	s.inputs[credIdxSystem].EchoMode = textinput.EchoPassword
	s.inputs[credIdxSystem].EchoCharacter = '*'
	s.inputs[credIdxSystem].CharLimit = 30

	// PDBADMIN Password
	s.inputs[credIdxPDBAdmin] = textinput.New()
	s.inputs[credIdxPDBAdmin].Placeholder = "PDB Admin password"
	s.inputs[credIdxPDBAdmin].EchoMode = textinput.EchoPassword
	s.inputs[credIdxPDBAdmin].EchoCharacter = '*'
	s.inputs[credIdxPDBAdmin].CharLimit = 30

	return s
}

// Init initializes the step
func (s *CredentialsStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""
	s.useCommonPassword = config.UseCommonPassword

	// Set values from config
	s.inputs[credIdxCommon].SetValue(config.CommonPassword)
	s.inputs[credIdxSys].SetValue(config.SysPassword)
	s.inputs[credIdxSystem].SetValue(config.SystemPassword)
	s.inputs[credIdxPDBAdmin].SetValue(config.PDBAdminPassword)

	// Reset focus
	for i := range s.inputs {
		s.inputs[i].Blur()
	}

	return nil
}

// Update handles messages
func (s *CredentialsStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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
			// Toggle common password mode
			if s.focusIndex == 0 {
				s.useCommonPassword = !s.useCommonPassword
				if s.useCommonPassword {
					s.focusIndex = 1
					s.inputs[credIdxCommon].Focus()
				} else {
					s.focusIndex = 1
					s.inputs[credIdxSys].Focus()
				}
				return s, wizard.StepStay, textinput.Blink
			}
		}
	}

	// Update the focused text input
	inputIdx := s.getInputIndex()
	if inputIdx >= 0 {
		var cmd tea.Cmd
		s.inputs[inputIdx], cmd = s.inputs[inputIdx].Update(msg)
		return s, wizard.StepStay, cmd
	}

	return s, wizard.StepStay, nil
}

func (s *CredentialsStep) getInputIndex() int {
	if s.focusIndex == 0 {
		return -1 // Toggle
	}

	if s.useCommonPassword {
		if s.focusIndex == 1 {
			return credIdxCommon
		}
		return -1
	}

	// Different passwords mode
	switch s.focusIndex {
	case 1:
		return credIdxSys
	case 2:
		return credIdxSystem
	case 3:
		if s.config.CreateAsContainerDB {
			return credIdxPDBAdmin
		}
	}
	return -1
}

func (s *CredentialsStep) getMaxFields() int {
	if s.useCommonPassword {
		return 1 // toggle + common password
	}
	if s.config.CreateAsContainerDB {
		return 3 // toggle + sys + system + pdbadmin
	}
	return 2 // toggle + sys + system
}

func (s *CredentialsStep) nextField() {
	inputIdx := s.getInputIndex()
	if inputIdx >= 0 {
		s.inputs[inputIdx].Blur()
	}

	maxFields := s.getMaxFields()
	s.focusIndex++
	if s.focusIndex > maxFields {
		s.focusIndex = 0
	}

	newInputIdx := s.getInputIndex()
	if newInputIdx >= 0 {
		s.inputs[newInputIdx].Focus()
	}
}

func (s *CredentialsStep) prevField() {
	inputIdx := s.getInputIndex()
	if inputIdx >= 0 {
		s.inputs[inputIdx].Blur()
	}

	maxFields := s.getMaxFields()
	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxFields
	}

	newInputIdx := s.getInputIndex()
	if newInputIdx >= 0 {
		s.inputs[newInputIdx].Focus()
	}
}

func (s *CredentialsStep) validate() bool {
	s.err = ""

	if s.useCommonPassword {
		pwd := s.inputs[credIdxCommon].Value()
		if pwd == "" {
			s.err = "Password is required"
			return false
		}
		if len(pwd) < 8 {
			s.err = "Password must be at least 8 characters"
			return false
		}
	} else {
		if s.inputs[credIdxSys].Value() == "" {
			s.err = "SYS password is required"
			return false
		}
		if s.inputs[credIdxSystem].Value() == "" {
			s.err = "SYSTEM password is required"
			return false
		}
		if s.config.CreateAsContainerDB && s.inputs[credIdxPDBAdmin].Value() == "" {
			s.err = "PDB Admin password is required"
			return false
		}

		// Check minimum length
		if len(s.inputs[credIdxSys].Value()) < 8 {
			s.err = "SYS password must be at least 8 characters"
			return false
		}
	}

	return true
}

// View renders the step
func (s *CredentialsStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure database credentials:") + "\n\n")

	// Password mode toggle
	checkbox := ui.UncheckedStyle.String()
	if s.useCommonPassword {
		checkbox = ui.CheckedStyle.String()
	}
	toggleStyle := ui.NormalItemStyle
	if s.focusIndex == 0 {
		toggleStyle = ui.SelectedItemStyle
	}
	b.WriteString(checkbox + " " + toggleStyle.Render("Use same password for all accounts") + "\n")
	b.WriteString(ui.SubtitleStyle.Render("    Press 'c' to toggle") + "\n\n")

	if s.useCommonPassword {
		b.WriteString(s.renderField("Password for all accounts (SYS, SYSTEM, PDBADMIN)", s.inputs[credIdxCommon], s.focusIndex == 1))
	} else {
		b.WriteString(s.renderField("SYS Password", s.inputs[credIdxSys], s.focusIndex == 1))
		b.WriteString(s.renderField("SYSTEM Password", s.inputs[credIdxSystem], s.focusIndex == 2))
		if s.config.CreateAsContainerDB {
			b.WriteString(s.renderField("PDB Admin Password", s.inputs[credIdxPDBAdmin], s.focusIndex == 3))
		}
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Password requirements: minimum 8 characters") + "\n")

	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *CredentialsStep) renderField(label string, input textinput.Model, focused bool) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if focused {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	) + "\n\n"
}

// Title returns the step title
func (s *CredentialsStep) Title() string {
	return "Database Credentials"
}

// Apply applies the step's changes to the config
func (s *CredentialsStep) Apply(config *model.DBConfig) {
	config.UseCommonPassword = s.useCommonPassword

	if s.useCommonPassword {
		config.CommonPassword = s.inputs[credIdxCommon].Value()
		config.SysPassword = config.CommonPassword
		config.SystemPassword = config.CommonPassword
		config.PDBAdminPassword = config.CommonPassword
	} else {
		config.SysPassword = s.inputs[credIdxSys].Value()
		config.SystemPassword = s.inputs[credIdxSystem].Value()
		if config.CreateAsContainerDB {
			config.PDBAdminPassword = s.inputs[credIdxPDBAdmin].Value()
		}
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *CredentialsStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
