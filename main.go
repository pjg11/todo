package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var p *tea.Program

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.desc }

type model struct {
	list   list.Model
	active int
}

func (m model) Init() tea.Cmd {
	return nil
}

// Custom function to switch between the different categories.
func (m model) ChangeFilter(s string) tea.Cmd {
	go func() {
		if m.list.IsFiltered() {
			p.Send(tea.KeyMsg(tea.Key{Type: tea.KeyEsc}))
		}
		p.Send(tea.KeyMsg(tea.Key{Type: tea.KeyRunes, Runes: []rune{'/'}}))
		p.Send(tea.KeyMsg(tea.Key{Type: tea.KeyRunes, Runes: []rune(s)}))
		p.Send(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))
	}()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "tab" {
			switch m.active {
			case 0:
				m.ChangeFilter("today")
			case 1:
				m.ChangeFilter("todo")
			case 2:
				m.ChangeFilter("waiting")
			case 3:
				m.ChangeFilter("done")
			case 4:
				m.ChangeFilter("idea")
			case 5:
				m.ChangeFilter("archive")
			}
			m.active = (m.active + 1) % 6
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{}

	f, err := os.Open("list.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		line := scanner.Text()
		desc := " "
		if len(line) < 2 {
			continue
		}
		if line[1] == ' ' {
			switch line[0] {
			case '-':
				desc = "todo"
			case '*':
				desc = "today"
			case '/':
				desc = "waiting"
			case 'X':
				desc = "done"
			case '!':
				desc = "idea"
			case '<':
				desc = "archive"
			default:
				continue
			}
			items = append(items, item{title: line[2:], desc: desc})
		}
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "TODO"
	p = tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// TODO add help for switching between sections
// TODO Add feature to open editor from the app
// TODO Add help for editor
