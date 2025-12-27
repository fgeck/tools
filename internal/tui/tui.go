package tui

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/service"
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("46")) // Bright green
	itemStyle  = lipgloss.NewStyle().PaddingLeft(4)
	helpStyle  = lipgloss.NewStyle().PaddingLeft(4).PaddingTop(1).Foreground(lipgloss.Color("240"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	baseStyle  = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("34")) // Green
)

type tableRow struct {
	toolName    string
	description string // Example description
	command     string // The actual command to execute
}

type mode int

const (
	modeList mode = iota
	modeAdd
	modeEdit
	modeDelete
)

type model struct {
	table       table.Model
	tableRows   []tableRow
	service     service.ExampleService
	mode        mode
	err         error
	quitting    bool
	selectedCmd string // Command to output when exiting

	// Add/Edit mode fields
	toolNameInput textinput.Model
	descInput     textinput.Model
	cmdInput      textinput.Model
	focusIndex    int
	inputs        []textinput.Model

	// Edit mode specific
	originalCmd string // Original command being edited
}

type examplesLoadedMsg struct {
	examples []dto.ExampleResponse
}

type errorMsg struct {
	err error
}

func loadExamples(svc service.ExampleService) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := svc.ListExamples(ctx)
		if err != nil {
			return errorMsg{err}
		}
		return examplesLoadedMsg{examples: resp.Examples}
	}
}

func NewModel(svc service.ExampleService) model {
	columns := []table.Column{
		{Title: "Tool", Width: 15},
		{Title: "Description", Width: 40},
		{Title: "Command", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("34")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("46")) // Bright green
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("0")).  // Black text
		Background(lipgloss.Color("46")). // Bright green background
		Bold(false)
	t.SetStyles(s)

	// Initialize text inputs for add mode
	toolNameInput := textinput.New()
	toolNameInput.Placeholder = "Tool name (e.g., lsof)"
	toolNameInput.Focus()
	toolNameInput.CharLimit = 50
	toolNameInput.Width = 50

	descInput := textinput.New()
	descInput.Placeholder = "Description (e.g., list all ports at port 54321)"
	descInput.CharLimit = 200
	descInput.Width = 50

	cmdInput := textinput.New()
	cmdInput.Placeholder = "Command (e.g., lsof -i :54321)"
	cmdInput.CharLimit = 200
	cmdInput.Width = 50

	m := model{
		table:         t,
		service:       svc,
		mode:          modeList,
		toolNameInput: toolNameInput,
		descInput:     descInput,
		cmdInput:      cmdInput,
		inputs:        []textinput.Model{toolNameInput, descInput, cmdInput},
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadExamples(m.service), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 10)
		m.table.SetWidth(msg.Width)
		return m, nil

	case examplesLoadedMsg:
		rows := []table.Row{}
		m.tableRows = []tableRow{}
		for _, example := range msg.examples {
			rows = append(rows, table.Row{
				example.ToolName,
				example.Description,
				example.Command,
			})
			m.tableRows = append(m.tableRows, tableRow{
				toolName:    example.ToolName,
				description: example.Description,
				command:     example.Command,
			})
		}
		m.table.SetRows(rows)
		return m, nil

	case errorMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.handleListKeys(msg)
		case modeAdd:
			return m.handleAddKeys(msg)
		case modeEdit:
			return m.handleEditKeys(msg)
		case modeDelete:
			return m.handleDeleteKeys(msg)
		}
	}

	// Update table
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "q":
		m.quitting = true
		return m, tea.Quit

	case "a":
		m.mode = modeAdd
		m.focusIndex = 0
		m.inputs[0].Focus()
		return m, textinput.Blink

	case "e", "edit":
		if len(m.tableRows) > 0 {
			cursor := m.table.Cursor()
			if cursor >= 0 && cursor < len(m.tableRows) {
				row := m.tableRows[cursor]
				m.mode = modeEdit
				m.originalCmd = row.command
				// Pre-fill inputs with current values
				m.inputs[0].SetValue(row.toolName)
				m.inputs[1].SetValue(row.description)
				m.inputs[2].SetValue(row.command)
				m.focusIndex = 0
				m.inputs[0].Focus()
				return m, textinput.Blink
			}
		}

	case "d", "delete":
		if len(m.tableRows) > 0 {
			m.mode = modeDelete
			return m, nil
		}

	case "enter":
		// Select the command and exit
		cursor := m.table.Cursor()
		if cursor >= 0 && cursor < len(m.tableRows) {
			m.selectedCmd = m.tableRows[cursor].command
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) handleAddKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = modeList
		m.resetInputs()
		return m, nil

	case "enter":
		// Submit on enter from any field
		return m.submitAdd()

	case "tab", "shift+tab", "up", "down":
		s := msg.String()

		// Navigation
		switch s {
		case "up", "shift+tab":
			m.focusIndex--
		case "down", "tab":
			m.focusIndex++
		}

		if m.focusIndex > len(m.inputs)-1 {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i < len(m.inputs); i++ {
			if i == m.focusIndex {
				cmds[i] = m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}

		return m, tea.Batch(cmds...)
	}

	// Update current input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) handleEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = modeList
		m.resetInputs()
		return m, nil

	case "enter":
		// Submit on enter from any field
		return m.submitEdit()

	case "tab", "shift+tab", "up", "down":
		s := msg.String()

		// Navigation
		switch s {
		case "up", "shift+tab":
			m.focusIndex--
		case "down", "tab":
			m.focusIndex++
		}

		if m.focusIndex > len(m.inputs)-1 {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}

		cmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i < len(m.inputs); i++ {
			if i == m.focusIndex {
				cmds[i] = m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}

		return m, tea.Batch(cmds...)
	}

	// Update current input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) handleDeleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "n":
		m.mode = modeList
		return m, nil

	case "y", "enter":
		return m.submitDelete()
	}

	return m, nil
}

func (m *model) updateInputs(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Sync back to individual fields
	m.toolNameInput = m.inputs[0]
	m.descInput = m.inputs[1]
	m.cmdInput = m.inputs[2]

	return tea.Batch(cmds...)
}

func (m *model) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	m.toolNameInput.SetValue("")
	m.descInput.SetValue("")
	m.cmdInput.SetValue("")
	m.focusIndex = 0
}

func (m model) submitAdd() (tea.Model, tea.Cmd) {
	toolName := strings.TrimSpace(m.toolNameInput.Value())
	desc := strings.TrimSpace(m.descInput.Value())
	cmd := strings.TrimSpace(m.cmdInput.Value())

	if toolName == "" || desc == "" || cmd == "" {
		m.err = fmt.Errorf("tool name, description, and command are required")
		return m, nil
	}

	req := dto.CreateExampleRequest{
		Command:     cmd,
		ToolName:    toolName,
		Description: desc,
	}

	ctx := context.Background()
	_, err := m.service.CreateExample(ctx, req)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.mode = modeList
	m.resetInputs()
	m.err = nil
	return m, loadExamples(m.service)
}

func (m model) submitEdit() (tea.Model, tea.Cmd) {
	toolName := strings.TrimSpace(m.toolNameInput.Value())
	desc := strings.TrimSpace(m.descInput.Value())
	cmd := strings.TrimSpace(m.cmdInput.Value())

	if toolName == "" || desc == "" || cmd == "" {
		m.err = fmt.Errorf("tool name, description, and command are required")
		return m, nil
	}

	req := dto.UpdateExampleRequest{
		Command:        m.originalCmd,
		NewToolName:    toolName,
		NewDescription: desc,
		NewCommand:     cmd,
	}

	ctx := context.Background()
	_, err := m.service.UpdateExample(ctx, req)
	if err != nil {
		m.err = err
		return m, nil
	}

	m.mode = modeList
	m.resetInputs()
	m.err = nil
	return m, loadExamples(m.service)
}

func (m model) submitDelete() (tea.Model, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.tableRows) {
		return m, nil
	}

	row := m.tableRows[cursor]
	ctx := context.Background()
	// Delete the specific example by its command (primary key)
	err := m.service.DeleteExample(ctx, row.command)
	if err != nil {
		m.err = err
		m.mode = modeList
		return m, nil
	}

	m.mode = modeList
	m.err = nil
	return m, loadExamples(m.service)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	switch m.mode {
	case modeAdd:
		return m.addView()
	case modeEdit:
		return m.editView()
	case modeDelete:
		return m.deleteView()
	default:
		return m.listView()
	}
}

func (m model) listView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Tools - Command Bookmarks"))
	b.WriteString("\n\n")
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteString("\n")

	// Help
	help := helpStyle.Render("↑/↓: navigate • enter: select (copies to clipboard) • a: add • e: edit • d: delete • q/esc: quit")
	b.WriteString(help)

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

func (m model) addView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add New Example"))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Tool Name:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[0].View()))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Description:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[1].View()))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Command:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[2].View()))
	b.WriteString("\n\n")

	help := helpStyle.Render("tab/shift+tab: navigate • enter: submit • esc: cancel")
	b.WriteString(help)

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

func (m model) editView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Edit Example"))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Tool Name:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[0].View()))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Description:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[1].View()))
	b.WriteString("\n\n")

	b.WriteString(itemStyle.Render("Command:"))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(m.inputs[2].View()))
	b.WriteString("\n\n")

	help := helpStyle.Render("tab/shift+tab: navigate • enter: submit • esc: cancel")
	b.WriteString(help)

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

func (m model) deleteView() string {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.tableRows) {
		return ""
	}

	row := m.tableRows[cursor]
	var b strings.Builder
	b.WriteString(titleStyle.Render("Confirm Delete"))
	b.WriteString("\n\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("Delete example '%s' from tool '%s'?", row.description, row.toolName)))
	b.WriteString("\n")
	b.WriteString(itemStyle.Render(fmt.Sprintf("Command: %s", row.command)))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("y: yes • n/esc: no"))

	return b.String()
}

func Run(svc service.ExampleService) error {
	m := NewModel(svc)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Output the selected command if one was chosen
	if fm, ok := finalModel.(model); ok && fm.selectedCmd != "" {
		// Copy to clipboard using OSC 52 escape sequence
		copyToClipboard(fm.selectedCmd)

		// Print the command to stdout
		fmt.Println(fm.selectedCmd)
	}

	return nil
}

// copyToClipboard uses OSC 52 escape sequence to copy to clipboard
func copyToClipboard(text string) {
	// Base64 encode the text
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	// OSC 52 escape sequence: \033]52;c;base64\007
	fmt.Printf("\033]52;c;%s\007", encoded)
}
