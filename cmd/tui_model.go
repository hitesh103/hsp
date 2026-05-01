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

	type TUIModel struct {
		URLInput      textinput.Model
		CurlInput     textinput.Model
		ShowCurlInput bool
		Method        string
		Headers       map[string]string
		Body          string
		ActiveTab     int
		Err           error
		Quitting      bool
		Sending       bool
	}

	func initialModel() TUIModel {
		ti := textinput.New()
		ti.Placeholder = "https://api.example.com"
		ti.Focus()
		ti.CharLimit = 156
		ti.Width = 60
		ti.Prompt = " "

		ci := textinput.New()
		ci.Placeholder = "Paste cURL command here..."
		ci.Width = 60
		ci.Prompt = " > "

		return TUIModel{
			URLInput:  ti,
			CurlInput: ci,
			Method:    "GET",
			Headers:   make(map[string]string),
			ActiveTab: 0,
		}
	}

	func (m TUIModel) Init() tea.Cmd { return textinput.Blink }

	func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		var cmd tea.Cmd

		switch msg := msg.(type) {
		case editorFinishedMsg:
			if msg.err != nil {
				m.Err = msg.err
				return m, nil
			}
			m.Body = msg.content
			return m, nil

		case tea.KeyMsg:
			if m.ShowCurlInput {
				switch msg.String() {
				case "esc":
					m.ShowCurlInput = false
					m.Err = nil
					return m, nil
				case "enter":
					config, err := ParseCurl(m.CurlInput.Value())
					if err != nil {
						m.Err = err
					} else {
						m.URLInput.SetValue(config.URL)
						m.Method = config.Method
						m.Headers = config.Headers
						m.Body = config.Body
						m.ShowCurlInput = false
						m.Err = nil
					}
					return m, nil
				}
				m.CurlInput, cmd = m.CurlInput.Update(msg)
				return m, cmd
			}

			switch msg.String() {
			case "ctrl+c", "esc":
				m.Quitting = true
				return m, tea.Quit
			case "ctrl+s":
				m.Sending = true
				return m, tea.Quit
			case "tab":
				m.ActiveTab = (m.ActiveTab + 1) % 3
			case "shift+tab":
				m.ActiveTab = (m.ActiveTab - 1 + 3) % 3
			case "ctrl+i":
				m.ShowCurlInput = true
				m.CurlInput.Focus()
				m.CurlInput.SetValue("")
				return m, nil
			case "ctrl+e":
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}

				f, err := os.CreateTemp("", "hsp-body-*.txt")
				if err != nil {
					m.Err = err
					return m, nil
				}
				_, _ = f.WriteString(m.Body)
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

		m.URLInput, cmd = m.URLInput.Update(msg)
		return m, cmd
	}

	func (m TUIModel) View() string {
		if m.Quitting {
			return ""
		}

		if m.ShowCurlInput {
			header := titleStyle.Render(" HSP - Import cURL ")
			input := m.CurlInput.View()
			inputBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(0, 1).
				Width(62).
				Render(input)

			help := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("\n enter: import • esc: cancel ")

			content := lipgloss.JoinVertical(lipgloss.Left, header, "Paste your cURL command:", inputBox, help)

			if m.Err != nil {
				errorMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("\nError: " + m.Err.Error())
				content = lipgloss.JoinVertical(lipgloss.Left, content, errorMsg)
			}

			return lipgloss.Place(80, 20, lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().
					Border(lipgloss.DoubleBorder()).
					BorderForeground(lipgloss.Color("#7D56F4")).
					Padding(1).
					Render(content))
		}

		header := titleStyle.Render(" HSP - HTTP Superpowers ")

		method := methodStyle.Render(" " + m.Method + " ")
		url := m.URLInput.View()

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
			isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == m.ActiveTab

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
		switch m.ActiveTab {
		case 0:
			content = "Headers Section (Placeholder)\n\nKey: Value pairs will go here."
		case 1:
			content = "Params Section (Placeholder)\n\nQuery parameters will go here."
		case 2:
			if m.Body == "" {
				content = "Body is empty.\n\nPress ctrl+e to edit in your $EDITOR."
			} else {
				content = m.Body
			}
		}

		tabContent := tabWindowStyle.Render(content)

		help := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("\n tab: next tab • shift+tab: prev tab • ctrl+e: edit body • ctrl+i: import curl • ctrl+s: send • ctrl+c: quit ")
		return lipgloss.JoinVertical(lipgloss.Left, header, headerSection, "", row, tabContent, help) + "\n"
	}

