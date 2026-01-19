package wizard

import (
	"dbca_tui/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

// StepResult indicates the outcome of a step
type StepResult int

const (
	StepContinue     StepResult = iota // Continue to next step
	StepBack                           // Go back to previous step
	StepQuit                           // Quit the wizard
	StepStay                           // Stay on current step
	StepPrintAndQuit                   // Print command and quit
)

// Step is the interface that all wizard steps must implement
type Step interface {
	// Init initializes the step with the current config
	Init(config *model.DBConfig) tea.Cmd

	// Update handles messages and returns a result
	Update(msg tea.Msg) (Step, StepResult, tea.Cmd)

	// View renders the step
	View() string

	// Title returns the step title
	Title() string

	// Apply applies the step's changes to the config
	Apply(config *model.DBConfig)

	// ShouldSkip returns true if this step should be skipped
	ShouldSkip(config *model.DBConfig) bool
}

// BaseStep provides common functionality for steps
type BaseStep struct {
	config *model.DBConfig
}

// ShouldSkip by default returns false (don't skip)
func (b *BaseStep) ShouldSkip(config *model.DBConfig) bool {
	return false
}
