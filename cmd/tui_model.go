package cmd

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editorFinishedMsg struct {
	err     error
	content string
}

var (	titleStyle = lipgloss.NewStyle().
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

	tabStyle = lipgloss.NewStyle().
	                Border(lipgloss.NormalBorder(), true, true, false, true).
	                BorderForeground(lipgloss.Color("#874BFD")).
	                Padding(0, 2)

	activeTabStyle = tabStyle.Copy().
	                        Border(lipgloss.NormalBorder(), true, true, true, true).
	                        BorderForeground(lipgloss.Color("#874BFD")).
	                        Foreground(lipgloss.Color("#7D56F4")).
	                        Bold(true)

	tabWindowStyle = lipgloss.NewStyle().
	                        Border(lipgloss.NormalBorder(), true, true, true, true).
	                        BorderForeground(lipgloss.Color("#874BFD")).
	                        Padding(1).
	                        Width(70).
	                        Height(10)
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
	ti.CharLimit = 156
	ti.Width = 60
	ti.Prompt = " "

	return model{
	        urlInput:  ti,
	        method:    "GET",
	        headers:   make(map[string]string),
	        activeTab: 0,
	}
	}

	func (m model) Init() tea.Cmd { return textinput.Blink }

	func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		var cmd tea.Cmd

		switch msg := msg.(type) {
		case editorFinishedMsg:
			if msg.err != nil {
				m.err = msg.err
				return m, nil
			}
			m.body = msg.content
			return m, nil

		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.quitting = true
				return m, tea.Quit
			case "tab":
				m.activeTab = (m.activeTab + 1) % 3
			case "shift+tab":
				m.activeTab = (m.activeTab - 1 + 3) % 3
			case "ctrl+e":
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}

				f, err := os.CreateTemp("", "hsp-body-*.txt")
				if err != nil {
					m.err = err
					return m, nil
				}
				_, _ = f.WriteString(m.body)
				f.Close()

				c := exec.Command(editor, f.Name())
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					if err != nil {
						return editorFinishedMsg{err: err}
					}
					content, err := os.ReadFile(f.Name())
					defer os.Remove(f.Name())
					if err != nil {
						return editorFinishedMsg{err: err}
					}
					return editorFinishedMsg{content: string(content)}
				})
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

	// Tabs
	tabs := []string{"Headers", "Query Params", "Body"}
	var renderedTabs []string

	for i, t := range tabs {
	        var style lipgloss.Style
	        isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == m.activeTab

	        if isActive {
	                style = activeTabStyle.Copy()
	        } else {
	                style = tabStyle.Copy()
	        }

	        // Adjust borders for seamless look
	        if isFirst && !isActive {
	                style = style.Border(lipgloss.NormalBorder(), true, true, false, true)
	        } else if isLast && !isActive {
	                style = style.Border(lipgloss.NormalBorder(), true, true, false, true)
	        }

	        renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Content
	var content string
	switch m.activeTab {
	case 0:
	        content = "Headers Section (Placeholder)\n\nKey: Value pairs will go here."
	case 1:
	        content = "Params Section (Placeholder)\n\nQuery parameters will go here."
	case 2:
	        if m.body == "" {
	                content = "Body is empty.\n\nPress ctrl+e to edit in your $EDITOR."
	        } else {
	                content = m.body
	        }
	}

	tabContent := tabWindowStyle.Render(content)

	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("\n tab: next tab • shift+tab: prev tab • ctrl+e: edit body • ctrl+c: quit ")
	return lipgloss.JoinVertical(lipgloss.Left, header, headerSection, "", row, tabContent, help) + "\n"
	}


