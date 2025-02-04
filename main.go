package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)


func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

const (
	intialInput = 2
	maxInput = 6
	minInputs = 1
	helpHeight = 5
)

var (
	cursorStyle = lg.NewStyle().Foreground(lg.Color("255"))
	cursorLineStyle = lg.NewStyle().
		Background(lg.Color("88")).
		Foreground(lg.Color("255"))	
	placeholderStyle = lg.NewStyle().
		Foreground(lg.Color("238"))
	endOfBufferStyle = lg.NewStyle().
		Foreground(lg.Color("235"))
	focusedPlaceholderStyle = lg.NewStyle().
		Foreground(lg.Color("255"))
	focusedBorderStyle = lg.NewStyle().
		Border(lg.RoundedBorder()).
		BorderForeground(lg.Color("238"))
	blurredBorderStyle = lg.NewStyle().
		Border(lg.HiddenBorder())
)

type keymap = struct {
	next, prev, add, remove, quit key.Binding
}

func newTextarea(exaName string) textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = exaName
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()

	return t
}

type model struct {
	width int
	height int
	keymap keymap
	help help.Model
	inputs []textarea.Model
	focus int
}

func newModel() model {
	m := model {
		inputs: make([]textarea.Model, intialInput),
		help: help.New(),
		keymap: keymap {
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			add: key.NewBinding(
				key.WithKeys("ctrl+n"),
				key.WithHelp("ctrl+n", "add an editor"),
			),
			remove: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "remove an editor"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}

	for i := 0; i < intialInput; i++ {
		m.inputs[i] = newTextarea("X" + string(byte(65+i)))
	}
	m.inputs[m.focus].Focus()
	m.updateKeybindings()

	return m
}

func (m *model) updateKeybindings() {
	m.keymap.add.SetEnabled(len(m.inputs) < maxInput)
	m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m *model) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit

		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs) - 1 {
				m.focus = 0
			}

			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()			
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}

			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keymap.add):
			m.inputs = append(m.inputs, newTextarea("X" + string(byte(65+len(m.inputs)))))

		case key.Matches(msg, m.keymap.remove):
			tmp := m.inputs[:m.focus]
			m.inputs = append(tmp, m.inputs[m.focus+1:]...)
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.updateKeybindings()
	m.sizeInputs()

	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg) 
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.add,
		m.keymap.remove,
		m.keymap.quit,
	})

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	return lg.JoinHorizontal(lg.Top, views...) + "\n\n" + help
}
