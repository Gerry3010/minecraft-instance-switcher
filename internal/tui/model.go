package tui

import (
	"fmt"
	"strings"

	"github.com/Gerry3010/minecraft-instance-switcher/internal/instance"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateList state = iota
	stateDetail
	stateCreate
	stateConfirmDelete
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Back        key.Binding
	Quit        key.Binding
	Help        key.Binding
	Create      key.Binding
	Delete      key.Binding
	Refresh     key.Binding
	Restore     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Create, k.Delete, k.Refresh},
		{k.Restore, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/switch"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create instance"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete instance"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Restore: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restore default"),
	),
}

type instanceItem struct {
	instance.Instance
}

func (i instanceItem) FilterValue() string { return i.Name }
func (i instanceItem) Title() string       { return i.Name }
func (i instanceItem) Description() string {
	var status string
	if i.IsActive {
		status = "● ACTIVE"
	} else {
		status = "○ Inactive"
	}
	return fmt.Sprintf("%s | %d mods | %d configs | %d saves", 
		status, i.ModCount, i.ConfigCount, i.SaveCount)
}

type model struct {
	state          state
	manager        *instance.Manager
	list           list.Model
	help           help.Model
	textInput      textinput.Model
	instances      []instance.Instance
	selectedInstance *instance.Instance
	instanceInfo   *instance.InstanceInfo
	message        string
	err            error
	keys           keyMap
}

type refreshMsg struct{}
type switchMsg struct{ name string }
type createMsg struct{ name string }
type deleteMsg struct{ name string }
type restoreMsg struct{}

func initialModel() model {
	manager, err := instance.NewManager()
	
	// Initialize list
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Minecraft Instance Manager"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter instance name..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	// Initialize help
	h := help.New()

	m := model{
		state:     stateList,
		manager:   manager,
		list:      l,
		help:      h,
		textInput: ti,
		keys:      keys,
		err:       err,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		refreshInstances,
		textinput.Blink,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateDetail:
			return m.updateDetail(msg)
		case stateCreate:
			return m.updateCreate(msg)
		case stateConfirmDelete:
			return m.updateConfirmDelete(msg)
		}

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case refreshMsg:
		instances, err := m.manager.ListInstances()
		if err != nil {
			m.err = err
			return m, nil
		}
		
		m.instances = instances
		items := make([]list.Item, len(instances))
		for i, inst := range instances {
			items[i] = instanceItem{inst}
		}
		
		m.list.SetItems(items)
		m.err = nil
		return m, nil

	case switchMsg:
		err := m.manager.SwitchInstance(msg.name)
		if err != nil {
			m.err = err
		} else {
			m.message = fmt.Sprintf("Switched to instance: %s", msg.name)
		}
		return m, refreshInstances

	case createMsg:
		err := m.manager.CreateInstance(msg.name)
		if err != nil {
			m.err = err
		} else {
			m.message = fmt.Sprintf("Created instance: %s", msg.name)
			m.state = stateList
		}
		return m, refreshInstances

	case deleteMsg:
		err := m.manager.DeleteInstance(msg.name)
		if err != nil {
			m.err = err
		} else {
			m.message = fmt.Sprintf("Deleted instance: %s", msg.name)
			m.state = stateList
		}
		return m, refreshInstances

	case restoreMsg:
		err := m.manager.RestoreDefault()
		if err != nil {
			m.err = err
		} else {
			m.message = "Restored default minecraft directory"
		}
		return m, refreshInstances
	}

	// Update sub-models
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Enter):
		if len(m.instances) == 0 {
			return m, nil
		}
		
		selected := m.list.SelectedItem().(instanceItem)
		if selected.IsActive {
			// Show details if already active
			m.selectedInstance = &selected.Instance
			info, err := m.manager.GetInstanceInfo(selected.Name)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.instanceInfo = info
			m.state = stateDetail
		} else {
			// Switch to this instance
			return m, func() tea.Msg {
				return switchMsg{name: selected.Name}
			}
		}

	case key.Matches(msg, m.keys.Create):
		m.state = stateCreate
		m.textInput.SetValue("")
		m.textInput.Focus()

	case key.Matches(msg, m.keys.Delete):
		if len(m.instances) == 0 {
			return m, nil
		}
		
		selected := m.list.SelectedItem().(instanceItem)
		if selected.IsActive {
			m.err = fmt.Errorf("cannot delete active instance")
			return m, nil
		}
		
		m.selectedInstance = &selected.Instance
		m.state = stateConfirmDelete

	case key.Matches(msg, m.keys.Refresh):
		m.message = ""
		m.err = nil
		return m, refreshInstances

	case key.Matches(msg, m.keys.Restore):
		return m, func() tea.Msg {
			return restoreMsg{}
		}

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Quit):
		m.state = stateList
	}
	return m, nil
}

func (m model) updateCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Enter):
		name := strings.TrimSpace(m.textInput.Value())
		if name != "" {
			return m, func() tea.Msg {
				return createMsg{name: name}
			}
		}

	case key.Matches(msg, m.keys.Back):
		m.state = stateList
		m.textInput.Blur()
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, func() tea.Msg {
			return deleteMsg{name: m.selectedInstance.Name}
		}
	case "n", "N", "esc":
		m.state = stateList
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateList:
		return m.viewList()
	case stateDetail:
		return m.viewDetail()
	case stateCreate:
		return m.viewCreate()
	case stateConfirmDelete:
		return m.viewConfirmDelete()
	}
	return ""
}

func (m model) viewList() string {
	var content strings.Builder
	
	if m.err != nil {
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		content.WriteString("\n\n")
	}

	if m.message != "" {
		content.WriteString(successStyle.Render(m.message))
		content.WriteString("\n\n")
	}

	content.WriteString(m.list.View())
	content.WriteString("\n")
	content.WriteString(m.help.View(m.keys))

	return content.String()
}

func (m model) viewDetail() string {
	if m.selectedInstance == nil || m.instanceInfo == nil {
		return "No instance selected"
	}

	var content strings.Builder
	
	content.WriteString(titleStyle.Render(fmt.Sprintf("Instance Details: %s", m.selectedInstance.Name)))
	content.WriteString("\n\n")

	// Instance stats
	content.WriteString(subtitleStyle.Render("Statistics:"))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("• Mods: %d\n", m.selectedInstance.ModCount))
	content.WriteString(fmt.Sprintf("• Configs: %d\n", m.selectedInstance.ConfigCount))
	content.WriteString(fmt.Sprintf("• Saves: %d\n", m.selectedInstance.SaveCount))
	content.WriteString(fmt.Sprintf("• Status: %s\n", func() string {
		if m.selectedInstance.IsActive {
			return "Active"
		}
		return "Inactive"
	}()))
	content.WriteString("\n")

	// Mods list
	if len(m.instanceInfo.ModsDir) > 0 {
		content.WriteString(subtitleStyle.Render("Mods:"))
		content.WriteString("\n")
		for _, mod := range m.instanceInfo.ModsDir[:min(len(m.instanceInfo.ModsDir), 10)] {
			content.WriteString(fmt.Sprintf("• %s\n", mod))
		}
		if len(m.instanceInfo.ModsDir) > 10 {
			content.WriteString(fmt.Sprintf("... and %d more\n", len(m.instanceInfo.ModsDir)-10))
		}
		content.WriteString("\n")
	}

	// Saves list
	if len(m.instanceInfo.SavesDir) > 0 {
		content.WriteString(subtitleStyle.Render("Saves:"))
		content.WriteString("\n")
		for _, save := range m.instanceInfo.SavesDir {
			content.WriteString(fmt.Sprintf("• %s\n", save))
		}
		content.WriteString("\n")
	}

	content.WriteString(dimStyle.Render("Press ESC to go back"))

	return content.String()
}

func (m model) viewCreate() string {
	var content strings.Builder
	
	content.WriteString(titleStyle.Render("Create New Instance"))
	content.WriteString("\n\n")
	content.WriteString("Instance name:\n")
	content.WriteString(m.textInput.View())
	content.WriteString("\n\n")
	content.WriteString(dimStyle.Render("Press Enter to create, ESC to cancel"))

	return content.String()
}

func (m model) viewConfirmDelete() string {
	if m.selectedInstance == nil {
		return ""
	}

	var content strings.Builder
	
	content.WriteString(titleStyle.Render("Confirm Deletion"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Are you sure you want to delete instance '%s'?\n", m.selectedInstance.Name))
	content.WriteString("This action cannot be undone.\n\n")
	content.WriteString(errorStyle.Render("Press 'y' to confirm, 'n' to cancel"))

	return content.String()
}

func refreshInstances() tea.Msg {
	return refreshMsg{}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))
)