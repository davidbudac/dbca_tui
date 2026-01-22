package steps

import (
	"fmt"
	"strings"

	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeleteStep handles database deletion configuration
type DeleteStep struct {
	config      *model.DBConfig
	sidInput    textinput.Model
	sysPassword textinput.Model
	focusIndex  int
	forceDelete bool
	err         string
}

// NewDeleteStep creates a new delete step
func NewDeleteStep() *DeleteStep {
	s := &DeleteStep{}

	s.sidInput = textinput.New()
	s.sidInput.Placeholder = "orcl"
	s.sidInput.CharLimit = 12

	s.sysPassword = textinput.New()
	s.sysPassword.Placeholder = "SYS password"
	s.sysPassword.EchoMode = textinput.EchoPassword
	s.sysPassword.EchoCharacter = '*'
	s.sysPassword.CharLimit = 30

	return s
}

// Init initializes the step
func (s *DeleteStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""
	s.forceDelete = config.DeleteForce

	s.sidInput.SetValue(config.DeleteSID)
	s.sysPassword.SetValue(config.SysPassword)

	s.sidInput.Blur()
	s.sysPassword.Blur()
	s.sidInput.Focus()

	return textinput.Blink
}

// Update handles messages
func (s *DeleteStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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

		case "f", "F":
			if s.focusIndex == 2 {
				s.forceDelete = !s.forceDelete
			}
		}
	}

	// Update the focused text input
	var cmd tea.Cmd
	switch s.focusIndex {
	case 0:
		s.sidInput, cmd = s.sidInput.Update(msg)
	case 1:
		s.sysPassword, cmd = s.sysPassword.Update(msg)
	}

	return s, wizard.StepStay, cmd
}

func (s *DeleteStep) nextField() {
	switch s.focusIndex {
	case 0:
		s.sidInput.Blur()
		s.focusIndex = 1
		s.sysPassword.Focus()
	case 1:
		s.sysPassword.Blur()
		s.focusIndex = 2
	case 2:
		s.focusIndex = 0
		s.sidInput.Focus()
	}
}

func (s *DeleteStep) prevField() {
	switch s.focusIndex {
	case 0:
		s.sidInput.Blur()
		s.focusIndex = 2
	case 1:
		s.sysPassword.Blur()
		s.focusIndex = 0
		s.sidInput.Focus()
	case 2:
		s.focusIndex = 1
		s.sysPassword.Focus()
	}
}

func (s *DeleteStep) validate() bool {
	s.err = ""

	sid := strings.TrimSpace(s.sidInput.Value())
	if sid == "" {
		s.err = "Database SID is required"
		return false
	}
	if len(sid) > 12 {
		s.err = "SID must be 12 characters or less"
		return false
	}

	pwd := s.sysPassword.Value()
	if pwd == "" {
		s.err = "SYS password is required for deletion"
		return false
	}

	return true
}

// View renders the step
func (s *DeleteStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure database deletion:") + "\n\n")

	// Warning
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555")).
		Bold(true)
	b.WriteString(warningStyle.Render("WARNING: This will generate a command to permanently delete the database!") + "\n\n")

	// SID input
	b.WriteString(s.renderField("Database SID to delete", s.sidInput, 0) + "\n")

	// SYS Password
	b.WriteString(s.renderField("SYS Password", s.sysPassword, 1) + "\n")

	// Force delete toggle
	checkbox := ui.UncheckedStyle.String()
	if s.forceDelete {
		checkbox = ui.CheckedStyle.String()
	}
	forceStyle := ui.NormalItemStyle
	if s.focusIndex == 2 {
		forceStyle = ui.SelectedItemStyle
	}
	b.WriteString(fmt.Sprintf("\n%s %s\n", checkbox, forceStyle.Render("Force delete (abort running database)")))
	b.WriteString(ui.SubtitleStyle.Render("    Press 'f' to toggle") + "\n")

	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *DeleteStep) renderField(label string, input textinput.Model, index int) string {
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
func (s *DeleteStep) Title() string {
	return "Delete Database"
}

// Apply applies the step's changes to the config
func (s *DeleteStep) Apply(config *model.DBConfig) {
	config.DeleteSID = strings.TrimSpace(s.sidInput.Value())
	config.SysPassword = s.sysPassword.Value()
	config.DeleteForce = s.forceDelete
}

// ShouldSkip returns whether this step should be skipped
func (s *DeleteStep) ShouldSkip(config *model.DBConfig) bool {
	return config.Operation != model.OperationDelete
}
