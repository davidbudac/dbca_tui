package steps

import (
	"fmt"
	"strconv"
	"strings"

	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"
	"dbca_tui/internal/wizard"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RecoveryStep handles Fast Recovery Area configuration
type RecoveryStep struct {
	config         *model.DBConfig
	inputs         []textinput.Model
	focusIndex     int
	enableFRA      bool
	enableArchive  bool
	err            string
}

const (
	recIdxFRADest = iota
	recIdxFRASize
)

// NewRecoveryStep creates a new recovery step
func NewRecoveryStep() *RecoveryStep {
	s := &RecoveryStep{
		inputs: make([]textinput.Model, 2),
	}

	// FRA Destination
	s.inputs[recIdxFRADest] = textinput.New()
	s.inputs[recIdxFRADest].Placeholder = "/u01/app/oracle/fast_recovery_area"
	s.inputs[recIdxFRADest].CharLimit = 256

	// FRA Size (MB)
	s.inputs[recIdxFRASize] = textinput.New()
	s.inputs[recIdxFRASize].Placeholder = "10240"
	s.inputs[recIdxFRASize].CharLimit = 10

	return s
}

// Init initializes the step
func (s *RecoveryStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""
	s.enableFRA = config.EnableFRA
	s.enableArchive = config.EnableArchiveLog

	// Set values from config
	s.inputs[recIdxFRADest].SetValue(config.FRADestination)
	s.inputs[recIdxFRASize].SetValue(strconv.Itoa(config.FRASize))

	// Reset focus
	for i := range s.inputs {
		s.inputs[i].Blur()
	}

	if s.enableFRA {
		s.inputs[0].Focus()
	}

	return textinput.Blink
}

// Update handles messages
func (s *RecoveryStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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
			// Toggle FRA
			if s.focusIndex == 0 || !s.enableFRA {
				s.enableFRA = !s.enableFRA
				if s.enableFRA {
					s.focusIndex = 1
					s.inputs[recIdxFRADest].Focus()
				} else {
					for i := range s.inputs {
						s.inputs[i].Blur()
					}
					s.focusIndex = 0
				}
			}

		case "a", "A":
			// Toggle archive log mode
			if s.enableFRA && (s.focusIndex == 3 || s.focusIndex == len(s.inputs)+1) {
				s.enableArchive = !s.enableArchive
			}
		}
	}

	// Update the focused text input
	if s.enableFRA && s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		inputIdx := s.focusIndex - 1
		var cmd tea.Cmd
		s.inputs[inputIdx], cmd = s.inputs[inputIdx].Update(msg)
		return s, wizard.StepStay, cmd
	}

	return s, wizard.StepStay, nil
}

func (s *RecoveryStep) nextField() {
	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Blur()
	}

	maxFields := 0 // FRA toggle
	if s.enableFRA {
		maxFields = len(s.inputs) + 1 // inputs + archive toggle
	}

	s.focusIndex++
	if s.focusIndex > maxFields {
		s.focusIndex = 0
	}

	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Focus()
	}
}

func (s *RecoveryStep) prevField() {
	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Blur()
	}

	maxFields := 0
	if s.enableFRA {
		maxFields = len(s.inputs) + 1
	}

	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxFields
	}

	if s.focusIndex > 0 && s.focusIndex <= len(s.inputs) {
		s.inputs[s.focusIndex-1].Focus()
	}
}

func (s *RecoveryStep) validate() bool {
	s.err = ""

	if s.enableFRA {
		fraDest := strings.TrimSpace(s.inputs[recIdxFRADest].Value())
		if fraDest == "" {
			s.err = "Fast Recovery Area location is required"
			return false
		}

		fraSize, err := strconv.Atoi(strings.TrimSpace(s.inputs[recIdxFRASize].Value()))
		if err != nil || fraSize < 1024 {
			s.err = "FRA size must be at least 1024 MB"
			return false
		}
	}

	return true
}

// View renders the step
func (s *RecoveryStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure Fast Recovery Area (FRA):") + "\n\n")

	// Enable FRA toggle
	checkbox := ui.UncheckedStyle.String()
	if s.enableFRA {
		checkbox = ui.CheckedStyle.String()
	}
	fraStyle := ui.NormalItemStyle
	if s.focusIndex == 0 {
		fraStyle = ui.SelectedItemStyle
	}
	b.WriteString(fmt.Sprintf("%s %s\n", checkbox, fraStyle.Render("Enable Fast Recovery Area")))
	b.WriteString(ui.SubtitleStyle.Render("    Press 'f' to toggle") + "\n\n")

	if s.enableFRA {
		// FRA Destination
		b.WriteString(s.renderField("FRA Location", s.inputs[recIdxFRADest], 1) + "\n")

		// FRA Size
		b.WriteString(s.renderField("FRA Size (MB)", s.inputs[recIdxFRASize], 2) + "\n")

		// Archive log toggle
		archiveCheckbox := ui.UncheckedStyle.String()
		if s.enableArchive {
			archiveCheckbox = ui.CheckedStyle.String()
		}
		archiveStyle := ui.NormalItemStyle
		if s.focusIndex == 3 {
			archiveStyle = ui.SelectedItemStyle
		}
		b.WriteString(fmt.Sprintf("\n%s %s\n", archiveCheckbox, archiveStyle.Render("Enable Archive Log Mode")))
		b.WriteString(ui.SubtitleStyle.Render("    Press 'a' to toggle") + "\n")
	}

	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *RecoveryStep) renderField(label string, input textinput.Model, fieldIndex int) string {
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
func (s *RecoveryStep) Title() string {
	return "Fast Recovery Area"
}

// Apply applies the step's changes to the config
func (s *RecoveryStep) Apply(config *model.DBConfig) {
	config.EnableFRA = s.enableFRA
	config.EnableArchiveLog = s.enableArchive

	if s.enableFRA {
		config.FRADestination = strings.TrimSpace(s.inputs[recIdxFRADest].Value())
		fraSize, _ := strconv.Atoi(strings.TrimSpace(s.inputs[recIdxFRASize].Value()))
		config.FRASize = fraSize
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *RecoveryStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
