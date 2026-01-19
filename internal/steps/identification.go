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

// IdentificationStep handles database identification
type IdentificationStep struct {
	config       *model.DBConfig
	inputs       []textinput.Model
	focusIndex   int
	createCDB    bool
	err          string
}

const (
	idxGlobalName = iota
	idxSID
	idxNumPDBs
	idxPDBName
)

// NewIdentificationStep creates a new identification step
func NewIdentificationStep() *IdentificationStep {
	s := &IdentificationStep{
		inputs: make([]textinput.Model, 4),
	}

	// Global Database Name
	s.inputs[idxGlobalName] = textinput.New()
	s.inputs[idxGlobalName].Placeholder = "orcl.example.com"
	s.inputs[idxGlobalName].CharLimit = 128

	// SID
	s.inputs[idxSID] = textinput.New()
	s.inputs[idxSID].Placeholder = "orcl"
	s.inputs[idxSID].CharLimit = 12

	// Number of PDBs
	s.inputs[idxNumPDBs] = textinput.New()
	s.inputs[idxNumPDBs].Placeholder = "1"
	s.inputs[idxNumPDBs].CharLimit = 3

	// PDB Name/Prefix
	s.inputs[idxPDBName] = textinput.New()
	s.inputs[idxPDBName].Placeholder = "orclpdb"
	s.inputs[idxPDBName].CharLimit = 30

	return s
}

// Init initializes the step
func (s *IdentificationStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.focusIndex = 0
	s.err = ""

	// Set values from config
	s.inputs[idxGlobalName].SetValue(config.GlobalDBName)
	s.inputs[idxSID].SetValue(config.SID)
	s.inputs[idxNumPDBs].SetValue(strconv.Itoa(config.NumberOfPDBs))
	s.inputs[idxPDBName].SetValue(config.PDBName)
	s.createCDB = config.CreateAsContainerDB

	// Focus first input
	for i := range s.inputs {
		s.inputs[i].Blur()
	}
	s.inputs[0].Focus()

	return textinput.Blink
}

// Update handles messages
func (s *IdentificationStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
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
			// Toggle CDB mode when not in text input
			if s.focusIndex >= len(s.inputs) {
				s.createCDB = !s.createCDB
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

func (s *IdentificationStep) nextField() {
	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Blur()
	}

	maxFields := len(s.inputs)
	if s.createCDB {
		maxFields++ // CDB toggle
	}

	s.focusIndex++
	if s.focusIndex > maxFields {
		s.focusIndex = 0
	}

	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *IdentificationStep) prevField() {
	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Blur()
	}

	maxFields := len(s.inputs)
	if s.createCDB {
		maxFields++
	}

	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxFields
	}

	if s.focusIndex < len(s.inputs) {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *IdentificationStep) validate() bool {
	s.err = ""

	// Validate Global DB Name
	globalName := strings.TrimSpace(s.inputs[idxGlobalName].Value())
	if globalName == "" {
		s.err = "Global Database Name is required"
		return false
	}

	// Validate SID
	sid := strings.TrimSpace(s.inputs[idxSID].Value())
	if sid == "" {
		s.err = "SID is required"
		return false
	}
	if len(sid) > 12 {
		s.err = "SID must be 12 characters or less"
		return false
	}

	// Validate PDB settings if CDB
	if s.createCDB {
		numPDBs, err := strconv.Atoi(strings.TrimSpace(s.inputs[idxNumPDBs].Value()))
		if err != nil || numPDBs < 0 || numPDBs > 252 {
			s.err = "Number of PDBs must be between 0 and 252"
			return false
		}

		if numPDBs > 0 {
			pdbName := strings.TrimSpace(s.inputs[idxPDBName].Value())
			if pdbName == "" {
				s.err = "PDB Name/Prefix is required when creating PDBs"
				return false
			}
		}
	}

	return true
}

// View renders the step
func (s *IdentificationStep) View() string {
	var b strings.Builder

	b.WriteString(ui.SubtitleStyle.Render("Configure database identification:") + "\n\n")

	// Global Database Name
	b.WriteString(s.renderField("Global Database Name", s.inputs[idxGlobalName], 0) + "\n")

	// SID
	b.WriteString(s.renderField("Oracle SID", s.inputs[idxSID], 1) + "\n")

	// CDB Toggle
	checkbox := ui.UncheckedStyle.String()
	if s.createCDB {
		checkbox = ui.CheckedStyle.String()
	}
	cdbStyle := ui.NormalItemStyle
	if s.focusIndex == len(s.inputs) {
		cdbStyle = ui.SelectedItemStyle
	}
	b.WriteString(fmt.Sprintf("\n%s %s\n", checkbox, cdbStyle.Render("Create as Container Database (CDB)")) + "\n")
	b.WriteString(ui.SubtitleStyle.Render("    Press 'c' to toggle") + "\n\n")

	// PDB settings (only if CDB enabled)
	if s.createCDB {
		b.WriteString(s.renderField("Number of PDBs", s.inputs[idxNumPDBs], 2) + "\n")
		b.WriteString(s.renderField("PDB Name/Prefix", s.inputs[idxPDBName], 3) + "\n")
	}

	// Error message
	if s.err != "" {
		b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
	}

	b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	return b.String()
}

func (s *IdentificationStep) renderField(label string, input textinput.Model, index int) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if s.focusIndex == index {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	)
}

// Title returns the step title
func (s *IdentificationStep) Title() string {
	return "Database Identification"
}

// Apply applies the step's changes to the config
func (s *IdentificationStep) Apply(config *model.DBConfig) {
	config.GlobalDBName = strings.TrimSpace(s.inputs[idxGlobalName].Value())
	config.SID = strings.TrimSpace(s.inputs[idxSID].Value())
	config.CreateAsContainerDB = s.createCDB

	if s.createCDB {
		numPDBs, _ := strconv.Atoi(strings.TrimSpace(s.inputs[idxNumPDBs].Value()))
		config.NumberOfPDBs = numPDBs
		config.PDBName = strings.TrimSpace(s.inputs[idxPDBName].Value())
		config.PDBPrefix = config.PDBName
	} else {
		config.NumberOfPDBs = 0
		config.PDBName = ""
		config.PDBPrefix = ""
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *IdentificationStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
