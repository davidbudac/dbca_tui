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

// StorageStep handles storage configuration
type StorageStep struct {
	config       *model.DBConfig
	storageList  ui.SelectList
	inputs       []textinput.Model
	focusIndex   int
	phase        int // 0 = selecting storage type, 1 = entering paths
	useOMF       bool
	err          string
}

const (
	stgIdxDatafile = iota
	stgIdxRedoLog
	stgIdxASMDiskGroup
)

// NewStorageStep creates a new storage step
func NewStorageStep() *StorageStep {
	items := []ui.SelectItem{
		{
			Title:       "File System",
			Description: "Store database files on a standard file system",
			Value:       string(model.StorageTypeFS),
		},
		{
			Title:       "Automatic Storage Management (ASM)",
			Description: "Store database files using Oracle ASM",
			Value:       string(model.StorageTypeASM),
		},
	}

	s := &StorageStep{
		storageList: ui.NewSelectList(items),
		inputs:      make([]textinput.Model, 3),
	}

	// Datafile destination
	s.inputs[stgIdxDatafile] = textinput.New()
	s.inputs[stgIdxDatafile].Placeholder = "/u01/app/oracle/oradata"
	s.inputs[stgIdxDatafile].CharLimit = 256

	// Redo log destination
	s.inputs[stgIdxRedoLog] = textinput.New()
	s.inputs[stgIdxRedoLog].Placeholder = "/u01/app/oracle/oradata"
	s.inputs[stgIdxRedoLog].CharLimit = 256

	// ASM Disk Group
	s.inputs[stgIdxASMDiskGroup] = textinput.New()
	s.inputs[stgIdxASMDiskGroup].Placeholder = "+DATA"
	s.inputs[stgIdxASMDiskGroup].CharLimit = 30

	return s
}

// Init initializes the step
func (s *StorageStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.phase = 0
	s.focusIndex = 0
	s.err = ""
	s.storageList.Reset()

	// Set cursor to current config value
	for i, item := range s.storageList.Items {
		if item.Value == string(config.StorageType) {
			s.storageList.Cursor = i
			break
		}
	}

	// Set input values
	s.inputs[stgIdxDatafile].SetValue(config.DatafileDestination)
	s.inputs[stgIdxRedoLog].SetValue(config.RedoLogDestination)
	s.inputs[stgIdxASMDiskGroup].SetValue(config.ASMDiskGroup)
	s.useOMF = config.UseOMF

	return nil
}

// Update handles messages
func (s *StorageStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.phase == 1 {
				s.phase = 0
				return s, wizard.StepStay, nil
			}
			return s, wizard.StepBack, nil
		}
	}

	if s.phase == 0 {
		return s.updateStorageSelection(msg)
	}
	return s.updatePathInput(msg)
}

func (s *StorageStep) updateStorageSelection(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			s.storageList.Update(msg)
			if s.storageList.IsSelected() {
				s.phase = 1
				s.focusIndex = 0
				for i := range s.inputs {
					s.inputs[i].Blur()
				}
				s.inputs[0].Focus()
				return s, wizard.StepStay, textinput.Blink
			}
		default:
			s.storageList.Update(msg)
		}
	}
	return s, wizard.StepStay, nil
}

func (s *StorageStep) updatePathInput(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	storageType := model.StorageType(s.storageList.GetSelectedValue())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			s.nextField(storageType)
			return s, wizard.StepStay, nil

		case "shift+tab", "up":
			s.prevField(storageType)
			return s, wizard.StepStay, nil

		case "enter":
			if s.validate(storageType) {
				return s, wizard.StepContinue, nil
			}
			return s, wizard.StepStay, nil

		case "o", "O":
			// Toggle OMF when in toggle position
			maxInputs := s.getMaxInputs(storageType)
			if s.focusIndex == maxInputs {
				s.useOMF = !s.useOMF
			}
		}
	}

	// Update the focused text input
	maxInputs := s.getMaxInputs(storageType)
	if s.focusIndex < maxInputs {
		var cmd tea.Cmd
		s.inputs[s.focusIndex], cmd = s.inputs[s.focusIndex].Update(msg)
		return s, wizard.StepStay, cmd
	}

	return s, wizard.StepStay, nil
}

func (s *StorageStep) getMaxInputs(storageType model.StorageType) int {
	if storageType == model.StorageTypeASM {
		return 1 // Just ASM disk group
	}
	return 2 // Datafile and redo log paths
}

func (s *StorageStep) nextField(storageType model.StorageType) {
	maxInputs := s.getMaxInputs(storageType)

	if s.focusIndex < maxInputs {
		s.inputs[s.focusIndex].Blur()
	}

	s.focusIndex++
	if s.focusIndex > maxInputs {
		s.focusIndex = 0
	}

	if s.focusIndex < maxInputs {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *StorageStep) prevField(storageType model.StorageType) {
	maxInputs := s.getMaxInputs(storageType)

	if s.focusIndex < maxInputs {
		s.inputs[s.focusIndex].Blur()
	}

	s.focusIndex--
	if s.focusIndex < 0 {
		s.focusIndex = maxInputs
	}

	if s.focusIndex < maxInputs {
		s.inputs[s.focusIndex].Focus()
	}
}

func (s *StorageStep) validate(storageType model.StorageType) bool {
	s.err = ""

	if storageType == model.StorageTypeASM {
		asmDG := strings.TrimSpace(s.inputs[stgIdxASMDiskGroup].Value())
		if asmDG == "" {
			s.err = "ASM Disk Group is required"
			return false
		}
	} else {
		datafile := strings.TrimSpace(s.inputs[stgIdxDatafile].Value())
		if datafile == "" {
			s.err = "Datafile destination is required"
			return false
		}
	}

	return true
}

// View renders the step
func (s *StorageStep) View() string {
	var b strings.Builder

	if s.phase == 0 {
		b.WriteString(ui.SubtitleStyle.Render("Select storage type:") + "\n\n")
		b.WriteString(s.storageList.View())
	} else {
		storageType := model.StorageType(s.storageList.GetSelectedValue())

		if storageType == model.StorageTypeASM {
			b.WriteString(ui.SubtitleStyle.Render("Configure ASM storage:") + "\n\n")
			b.WriteString(s.renderField("ASM Disk Group", s.inputs[stgIdxASMDiskGroup], 0) + "\n")
		} else {
			b.WriteString(ui.SubtitleStyle.Render("Configure file system storage:") + "\n\n")
			b.WriteString(s.renderField("Database Files Location", s.inputs[stgIdxDatafile], 0) + "\n")
			b.WriteString(s.renderField("Redo Log Files Location", s.inputs[stgIdxRedoLog], 1) + "\n")
		}

		// OMF Toggle
		maxInputs := s.getMaxInputs(storageType)
		checkbox := ui.UncheckedStyle.String()
		if s.useOMF {
			checkbox = ui.CheckedStyle.String()
		}
		omfStyle := ui.NormalItemStyle
		if s.focusIndex == maxInputs {
			omfStyle = ui.SelectedItemStyle
		}
		b.WriteString("\n" + checkbox + " " + omfStyle.Render("Use Oracle Managed Files (OMF)") + "\n")
		b.WriteString(ui.SubtitleStyle.Render("    Press 'o' to toggle") + "\n")

		if s.err != "" {
			b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
		}

		b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue, Esc to go back"))
	}

	return b.String()
}

func (s *StorageStep) renderField(label string, input textinput.Model, index int) string {
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
func (s *StorageStep) Title() string {
	return "Storage Configuration"
}

// Apply applies the step's changes to the config
func (s *StorageStep) Apply(config *model.DBConfig) {
	config.StorageType = model.StorageType(s.storageList.GetSelectedValue())
	config.UseOMF = s.useOMF

	if config.StorageType == model.StorageTypeASM {
		config.ASMDiskGroup = strings.TrimSpace(s.inputs[stgIdxASMDiskGroup].Value())
		config.DatafileDestination = config.ASMDiskGroup
		config.RedoLogDestination = config.ASMDiskGroup
	} else {
		config.DatafileDestination = strings.TrimSpace(s.inputs[stgIdxDatafile].Value())
		redoLog := strings.TrimSpace(s.inputs[stgIdxRedoLog].Value())
		if redoLog == "" {
			redoLog = config.DatafileDestination
		}
		config.RedoLogDestination = redoLog
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *StorageStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
