package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	urlInput    textinput.Model
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
	return model{
		urlInput: ti,
		method:   "GET",
		headers:  make(map[string]string),
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m model) View() string { return "TUI Initialized\n" }
