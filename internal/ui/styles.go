package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#FF6B35") // Oracle orange
	SecondaryColor = lipgloss.Color("#4A90D9") // Blue
	AccentColor    = lipgloss.Color("#00D4AA") // Teal
	TextColor      = lipgloss.Color("#FAFAFA")
	MutedColor     = lipgloss.Color("#888888")
	ErrorColor     = lipgloss.Color("#FF5555")
	SuccessColor   = lipgloss.Color("#55FF55")

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	// Subtitle/description style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginBottom(1)

	// Step indicator style
	StepIndicatorStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true)

	// Selected item style
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true)

	// Normal item style
	NormalItemStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// Cursor style
	CursorStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	// Input style
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(0, 1)

	// Focused input style
	FocusedInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(AccentColor).
				Padding(0, 1)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	// Box style for sections
	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Padding(1, 2)

	// Header style (for wizard header)
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 2).
			MarginBottom(1)

	// Code block style (for generated command)
	CodeBlockStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1A1A1A")).
			Foreground(lipgloss.Color("#E0E0E0")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Label style
	LabelStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Value style
	ValueStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// Checkbox styles
	CheckedStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			SetString("[x]")

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			SetString("[ ]")
)

// RenderHeader renders the wizard header
func RenderHeader(title string, step int, totalSteps int) string {
	stepText := StepIndicatorStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			"Step ",
			lipgloss.NewStyle().Bold(true).Render(string(rune('0'+step))),
			" of ",
			lipgloss.NewStyle().Render(string(rune('0'+totalSteps))),
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		HeaderStyle.Render("Oracle DBCA - Database Configuration Assistant"),
		TitleStyle.Render(title),
		stepText,
		"",
	)
}

// RenderHelp renders the navigation help
func RenderHelp() string {
	return HelpStyle.Render("↑/↓: Navigate • Enter: Select • Tab: Next field • Esc: Back • q: Quit")
}
