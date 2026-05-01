package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	methodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#EE6FF8")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
)

type model struct {	urlInput    textinput.Model
	method      string
	headers     map[string]string
	body        string
	activeTab   int
	err         error
	quitting    bool
	sending     bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60
	ti.Prompt = " "

	return model{
		urlInput: ti,
		method:   "GET",
		headers:  make(map[string]string),
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.urlInput, cmd = m.urlInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	header := titleStyle.Render(" HSP - HTTP Superpowers ")

	method := methodStyle.Render(" " + m.method + " ")
	url := m.urlInput.View()

	// Create a styled box for the URL input to make it look like a field
	urlBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1).
		Width(62).
		Render(url)

	headerSection := lipgloss.JoinHorizontal(lipgloss.Center, method, " ", urlBox)

	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(" ctrl+c: quit • esc: exit ")

	return lipgloss.JoinVertical(lipgloss.Left, header, headerSection, "", help) + "\n"
}

