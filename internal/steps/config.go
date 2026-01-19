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

// ConfigStep handles database configuration options (memory, charset, etc.)
type ConfigStep struct {
	config              *model.DBConfig
	memoryList          ui.SelectList
	charsetList         ui.SelectList
	connectionList      ui.SelectList
	memoryInput         textinput.Model
	phase               int // 0=memory type, 1=memory size, 2=charset, 3=connection mode
	enableSampleSchemas bool
	err                 string
}

// NewConfigStep creates a new config step
func NewConfigStep() *ConfigStep {
	memoryItems := []ui.SelectItem{
		{
			Title:       "Automatic Memory Management",
			Description: "Let Oracle automatically manage memory allocation (recommended)",
			Value:       "AUTO",
		},
		{
			Title:       "Automatic Shared Memory Management",
			Description: "Manually set total SGA, let Oracle manage PGA",
			Value:       "AUTO_SGA",
		},
		{
			Title:       "Manual Memory Management",
			Description: "Manually configure SGA and PGA sizes",
			Value:       "MANUAL",
		},
	}

	charsetItems := []ui.SelectItem{
		{
			Title:       "AL32UTF8 (Recommended)",
			Description: "Unicode UTF-8 Universal character set, supports all languages",
			Value:       "AL32UTF8",
		},
		{
			Title:       "UTF8",
			Description: "Unicode 3.0 UTF-8 Universal character set",
			Value:       "UTF8",
		},
		{
			Title:       "US7ASCII",
			Description: "US 7-bit ASCII character set",
			Value:       "US7ASCII",
		},
		{
			Title:       "WE8ISO8859P1",
			Description: "ISO 8859-1 West European character set",
			Value:       "WE8ISO8859P1",
		},
	}

	connectionItems := []ui.SelectItem{
		{
			Title:       "Dedicated Server Mode",
			Description: "Each client connection gets a dedicated server process",
			Value:       "DEDICATED",
		},
		{
			Title:       "Shared Server Mode",
			Description: "Multiple client connections share server processes",
			Value:       "SHARED",
		},
	}

	s := &ConfigStep{
		memoryList:     ui.NewSelectList(memoryItems),
		charsetList:    ui.NewSelectList(charsetItems),
		connectionList: ui.NewSelectList(connectionItems),
	}

	s.memoryInput = textinput.New()
	s.memoryInput.Placeholder = "2048"
	s.memoryInput.CharLimit = 10

	return s
}

// Init initializes the step
func (s *ConfigStep) Init(config *model.DBConfig) tea.Cmd {
	s.config = config
	s.phase = 0
	s.err = ""
	s.enableSampleSchemas = config.EnableSampleSchemas

	// Reset lists
	s.memoryList.Reset()
	s.charsetList.Reset()
	s.connectionList.Reset()

	// Set cursor positions based on config
	for i, item := range s.memoryList.Items {
		if item.Value == config.MemoryManagement {
			s.memoryList.Cursor = i
			break
		}
	}

	for i, item := range s.charsetList.Items {
		if item.Value == config.CharacterSet {
			s.charsetList.Cursor = i
			break
		}
	}

	for i, item := range s.connectionList.Items {
		if item.Value == config.ConnectionMode {
			s.connectionList.Cursor = i
			break
		}
	}

	s.memoryInput.SetValue(strconv.Itoa(config.TotalMemory))

	return nil
}

// Update handles messages
func (s *ConfigStep) Update(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if s.phase > 0 {
				s.phase--
				return s, wizard.StepStay, nil
			}
			return s, wizard.StepBack, nil
		}
	}

	switch s.phase {
	case 0:
		return s.updateMemoryType(msg)
	case 1:
		return s.updateMemorySize(msg)
	case 2:
		return s.updateCharset(msg)
	case 3:
		return s.updateConnectionMode(msg)
	}

	return s, wizard.StepStay, nil
}

func (s *ConfigStep) updateMemoryType(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			s.memoryList.Update(msg)
			if s.memoryList.IsSelected() {
				s.phase = 1
				s.memoryInput.Focus()
				return s, wizard.StepStay, textinput.Blink
			}
		default:
			s.memoryList.Update(msg)
		}
	}
	return s, wizard.StepStay, nil
}

func (s *ConfigStep) updateMemorySize(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			memSize, err := strconv.Atoi(strings.TrimSpace(s.memoryInput.Value()))
			if err != nil || memSize < 256 {
				s.err = "Memory size must be at least 256 MB"
				return s, wizard.StepStay, nil
			}
			s.err = ""
			s.phase = 2
			s.memoryInput.Blur()
			return s, wizard.StepStay, nil
		}
	}

	var cmd tea.Cmd
	s.memoryInput, cmd = s.memoryInput.Update(msg)
	return s, wizard.StepStay, cmd
}

func (s *ConfigStep) updateCharset(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			s.charsetList.Update(msg)
			if s.charsetList.IsSelected() {
				// In typical mode, skip connection mode selection
				if s.config.CreationMode == model.CreationModeTypical {
					return s, wizard.StepContinue, nil
				}
				s.phase = 3
				return s, wizard.StepStay, nil
			}
		default:
			s.charsetList.Update(msg)
		}
	}
	return s, wizard.StepStay, nil
}

func (s *ConfigStep) updateConnectionMode(msg tea.Msg) (wizard.Step, wizard.StepResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			s.connectionList.Update(msg)
			if s.connectionList.IsSelected() {
				return s, wizard.StepContinue, nil
			}
		case "s", "S":
			// Toggle sample schemas
			s.enableSampleSchemas = !s.enableSampleSchemas
		default:
			s.connectionList.Update(msg)
		}
	}
	return s, wizard.StepStay, nil
}

// View renders the step
func (s *ConfigStep) View() string {
	var b strings.Builder

	switch s.phase {
	case 0:
		b.WriteString(ui.SubtitleStyle.Render("Select memory management mode:") + "\n\n")
		b.WriteString(s.memoryList.View())

	case 1:
		memType := s.memoryList.GetSelectedItem()
		b.WriteString(ui.SubtitleStyle.Render(fmt.Sprintf("Memory Management: %s", memType.Title)) + "\n\n")
		b.WriteString(ui.LabelStyle.Render("Total Memory (MB)") + "\n")
		b.WriteString(ui.FocusedInputStyle.Render(s.memoryInput.View()) + "\n")
		b.WriteString(ui.SubtitleStyle.Render("    Recommended: At least 2048 MB") + "\n")

		if s.err != "" {
			b.WriteString("\n" + ui.ErrorStyle.Render(s.err) + "\n")
		}

		b.WriteString("\n" + ui.SubtitleStyle.Render("Press Enter to continue"))

	case 2:
		b.WriteString(ui.SubtitleStyle.Render("Select database character set:") + "\n\n")
		b.WriteString(s.charsetList.View())

	case 3:
		b.WriteString(ui.SubtitleStyle.Render("Select connection mode:") + "\n\n")
		b.WriteString(s.connectionList.View())

		// Sample schemas toggle
		checkbox := ui.UncheckedStyle.String()
		if s.enableSampleSchemas {
			checkbox = ui.CheckedStyle.String()
		}
		b.WriteString("\n\n" + checkbox + " " + ui.NormalItemStyle.Render("Install sample schemas (HR, OE, etc.)") + "\n")
		b.WriteString(ui.SubtitleStyle.Render("    Press 's' to toggle") + "\n")
	}

	return b.String()
}

func (s *ConfigStep) renderField(label string, input textinput.Model, focused bool) string {
	labelStyle := ui.LabelStyle
	inputStyle := ui.InputStyle

	if focused {
		inputStyle = ui.FocusedInputStyle
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(label),
		inputStyle.Render(input.View()),
	)
}

// Title returns the step title
func (s *ConfigStep) Title() string {
	return "Configuration Options"
}

// Apply applies the step's changes to the config
func (s *ConfigStep) Apply(config *model.DBConfig) {
	config.MemoryManagement = s.memoryList.GetSelectedValue()

	memSize, _ := strconv.Atoi(strings.TrimSpace(s.memoryInput.Value()))
	config.TotalMemory = memSize

	config.CharacterSet = s.charsetList.GetSelectedValue()
	config.NationalCharacterSet = "AL16UTF16" // Standard default

	if s.config.CreationMode == model.CreationModeAdvanced {
		config.ConnectionMode = s.connectionList.GetSelectedValue()
		config.EnableSampleSchemas = s.enableSampleSchemas
	} else {
		config.ConnectionMode = "DEDICATED"
		config.EnableSampleSchemas = false
	}
}

// ShouldSkip returns whether this step should be skipped
func (s *ConfigStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
