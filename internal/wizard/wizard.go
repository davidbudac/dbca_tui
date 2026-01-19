package wizard

import (
	"dbca_tui/internal/model"
	"dbca_tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Wizard is the main wizard model
type Wizard struct {
	steps        []Step
	currentStep  int
	config       *model.DBConfig
	width        int
	height       int
	quitting     bool
	completed    bool
	printCommand bool
}

// NewWizard creates a new wizard with the given steps
func NewWizard(steps []Step) *Wizard {
	return &Wizard{
		steps:       steps,
		currentStep: 0,
		config:      model.NewDBConfig(),
	}
}

// Init initializes the wizard
func (w *Wizard) Init() tea.Cmd {
	// Skip to first non-skippable step and initialize it
	w.skipToValidStep()
	if w.currentStep < len(w.steps) {
		return w.steps[w.currentStep].Init(w.config)
	}
	return nil
}

// Update handles messages
func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		return w, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			w.quitting = true
			return w, tea.Quit
		case "q":
			if w.completed {
				w.quitting = true
				return w, tea.Quit
			}
		}
	}

	if w.currentStep >= len(w.steps) || w.quitting {
		return w, nil
	}

	// Update current step
	step, result, cmd := w.steps[w.currentStep].Update(msg)
	w.steps[w.currentStep] = step

	switch result {
	case StepContinue:
		// Apply changes from current step
		w.steps[w.currentStep].Apply(w.config)

		// Move to next step
		w.currentStep++
		w.skipToValidStep()

		if w.currentStep >= len(w.steps) {
			w.completed = true
			return w, cmd
		}

		// Initialize next step
		initCmd := w.steps[w.currentStep].Init(w.config)
		return w, tea.Batch(cmd, initCmd)

	case StepBack:
		// Go back to previous step
		w.currentStep--
		w.skipBackToValidStep()

		if w.currentStep < 0 {
			w.currentStep = 0
		}

		// Re-initialize previous step
		initCmd := w.steps[w.currentStep].Init(w.config)
		return w, tea.Batch(cmd, initCmd)

	case StepQuit:
		w.quitting = true
		return w, tea.Quit

	case StepPrintAndQuit:
		w.quitting = true
		w.printCommand = true
		return w, tea.Quit
	}

	return w, cmd
}

// skipToValidStep skips forward to the next step that shouldn't be skipped
func (w *Wizard) skipToValidStep() {
	for w.currentStep < len(w.steps) && w.steps[w.currentStep].ShouldSkip(w.config) {
		w.currentStep++
	}
}

// skipBackToValidStep skips backward to the previous step that shouldn't be skipped
func (w *Wizard) skipBackToValidStep() {
	for w.currentStep > 0 && w.steps[w.currentStep].ShouldSkip(w.config) {
		w.currentStep--
	}
}

// View renders the wizard
func (w *Wizard) View() string {
	if w.quitting {
		return "Goodbye!\n"
	}

	if w.currentStep >= len(w.steps) {
		return "Wizard complete!\n"
	}

	// Count non-skipped steps for display
	totalSteps := 0
	currentDisplayStep := 0
	for i, step := range w.steps {
		if !step.ShouldSkip(w.config) {
			totalSteps++
			if i < w.currentStep {
				currentDisplayStep++
			} else if i == w.currentStep {
				currentDisplayStep++
			}
		}
	}

	step := w.steps[w.currentStep]

	// Build the view
	header := ui.RenderHeader(step.Title(), currentDisplayStep, totalSteps)
	content := step.View()
	help := ui.RenderHelp()

	// Calculate available height for content
	headerHeight := lipgloss.Height(header)
	helpHeight := lipgloss.Height(help)
	contentHeight := w.height - headerHeight - helpHeight - 2

	if contentHeight < 10 {
		contentHeight = 10
	}

	// Style the content area
	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		MaxHeight(contentHeight)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		contentStyle.Render(content),
		help,
	)
}

// GetConfig returns the current configuration
func (w *Wizard) GetConfig() *model.DBConfig {
	return w.config
}

// IsCompleted returns true if the wizard is completed
func (w *Wizard) IsCompleted() bool {
	return w.completed
}

// ShouldPrintCommand returns true if the command should be printed on exit
func (w *Wizard) ShouldPrintCommand() bool {
	return w.printCommand
}

// SetPrintCommand sets whether to print the command on exit
func (w *Wizard) SetPrintCommand(print bool) {
	w.printCommand = print
}
