package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectItem represents an item in a selection list
type SelectItem struct {
	Title       string
	Description string
	Value       string
}

// SelectList is a simple selection list component
type SelectList struct {
	Items    []SelectItem
	Cursor   int
	Selected int
}

// NewSelectList creates a new selection list
func NewSelectList(items []SelectItem) SelectList {
	return SelectList{
		Items:    items,
		Cursor:   0,
		Selected: -1,
	}
}

// Update handles input for the selection list
func (s *SelectList) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.Cursor > 0 {
				s.Cursor--
			}
		case "down", "j":
			if s.Cursor < len(s.Items)-1 {
				s.Cursor++
			}
		case "enter", " ":
			s.Selected = s.Cursor
		}
	}
}

// View renders the selection list
func (s SelectList) View() string {
	var b strings.Builder

	for i, item := range s.Items {
		cursor := "  "
		style := NormalItemStyle

		if i == s.Cursor {
			cursor = CursorStyle.Render("> ")
			style = SelectedItemStyle
		}

		title := style.Render(item.Title)
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, title))

		if item.Description != "" {
			desc := SubtitleStyle.Render("    " + item.Description)
			b.WriteString(desc + "\n")
		}
	}

	return b.String()
}

// IsSelected returns true if an item has been selected
func (s SelectList) IsSelected() bool {
	return s.Selected >= 0
}

// GetSelectedValue returns the value of the selected item
func (s SelectList) GetSelectedValue() string {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		return s.Items[s.Selected].Value
	}
	return ""
}

// GetSelectedItem returns the selected item
func (s SelectList) GetSelectedItem() *SelectItem {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		return &s.Items[s.Selected]
	}
	return nil
}

// Reset resets the selection
func (s *SelectList) Reset() {
	s.Selected = -1
}

// FormField represents a form input field
type FormField struct {
	Label       string
	Input       textinput.Model
	Placeholder string
	Required    bool
	Validator   func(string) error
	Error       string
}

// NewFormField creates a new form field
func NewFormField(label, placeholder string, required bool) FormField {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256

	return FormField{
		Label:       label,
		Input:       ti,
		Placeholder: placeholder,
		Required:    required,
	}
}

// NewPasswordField creates a new password field
func NewPasswordField(label, placeholder string, required bool) FormField {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '*'
	ti.CharLimit = 256

	return FormField{
		Label:       label,
		Input:       ti,
		Placeholder: placeholder,
		Required:    required,
	}
}

// Form is a simple form component
type Form struct {
	Fields       []FormField
	FocusedField int
}

// NewForm creates a new form
func NewForm(fields []FormField) Form {
	f := Form{
		Fields:       fields,
		FocusedField: 0,
	}

	if len(f.Fields) > 0 {
		f.Fields[0].Input.Focus()
	}

	return f
}

// Update handles input for the form
func (f *Form) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.NextField()
		case "shift+tab", "up":
			f.PrevField()
		}
	}

	// Update the focused field
	if f.FocusedField < len(f.Fields) {
		var cmd tea.Cmd
		f.Fields[f.FocusedField].Input, cmd = f.Fields[f.FocusedField].Input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// NextField moves to the next field
func (f *Form) NextField() {
	if f.FocusedField < len(f.Fields)-1 {
		f.Fields[f.FocusedField].Input.Blur()
		f.FocusedField++
		f.Fields[f.FocusedField].Input.Focus()
	}
}

// PrevField moves to the previous field
func (f *Form) PrevField() {
	if f.FocusedField > 0 {
		f.Fields[f.FocusedField].Input.Blur()
		f.FocusedField--
		f.Fields[f.FocusedField].Input.Focus()
	}
}

// View renders the form
func (f Form) View() string {
	var b strings.Builder

	for i, field := range f.Fields {
		// Label
		label := LabelStyle.Render(field.Label)
		if field.Required {
			label += ErrorStyle.Render(" *")
		}
		b.WriteString(label + "\n")

		// Input field
		inputStyle := InputStyle
		if i == f.FocusedField {
			inputStyle = FocusedInputStyle
		}
		b.WriteString(inputStyle.Render(field.Input.View()) + "\n")

		// Error message
		if field.Error != "" {
			b.WriteString(ErrorStyle.Render(field.Error) + "\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}

// GetValue returns the value of a field by index
func (f Form) GetValue(index int) string {
	if index >= 0 && index < len(f.Fields) {
		return f.Fields[index].Input.Value()
	}
	return ""
}

// SetValue sets the value of a field by index
func (f *Form) SetValue(index int, value string) {
	if index >= 0 && index < len(f.Fields) {
		f.Fields[index].Input.SetValue(value)
	}
}

// Validate validates all form fields
func (f *Form) Validate() bool {
	valid := true
	for i := range f.Fields {
		f.Fields[i].Error = ""

		if f.Fields[i].Required && f.Fields[i].Input.Value() == "" {
			f.Fields[i].Error = "This field is required"
			valid = false
		} else if f.Fields[i].Validator != nil {
			if err := f.Fields[i].Validator(f.Fields[i].Input.Value()); err != nil {
				f.Fields[i].Error = err.Error()
				valid = false
			}
		}
	}
	return valid
}

// Toggle represents a boolean toggle
type Toggle struct {
	Label   string
	Enabled bool
	Focused bool
}

// NewToggle creates a new toggle
func NewToggle(label string, enabled bool) Toggle {
	return Toggle{
		Label:   label,
		Enabled: enabled,
	}
}

// Update handles input for the toggle
func (t *Toggle) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			t.Enabled = !t.Enabled
		}
	}
}

// View renders the toggle
func (t Toggle) View() string {
	checkbox := UncheckedStyle.String()
	if t.Enabled {
		checkbox = CheckedStyle.String()
	}

	style := NormalItemStyle
	if t.Focused {
		style = SelectedItemStyle
	}

	return fmt.Sprintf("%s %s", checkbox, style.Render(t.Label))
}

// RenderKeyValue renders a key-value pair
func RenderKeyValue(key, value string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left,
		LabelStyle.Render(key+": "),
		ValueStyle.Render(value),
	)
}

// RenderSection renders a section with a title and content
func RenderSection(title, content string) string {
	return BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			TitleStyle.Render(title),
			content,
		),
	)
}
