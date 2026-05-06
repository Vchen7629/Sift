package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)


func main() {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error launching: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices  []string
	cursor 	 int
	selected map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{
		choices: []string{"Rape robin", "Rape Yae Miko", "Rape ahri"},

		selected: make(map[int]struct{}),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	// check if its a key press
	case tea.KeyPressMsg:
		// check what key was pressed
		switch msg.String() {
		
		case "ctrl+c", "q":
			return nil, tea.Quit
		
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j": 
			if m.cursor < len(m.choices) - 1 {
				m.cursor++
			}
		
		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}		
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Which girl should we knock up and rape\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		choiceSelected := " "
		_, ok := m.selected[i]
		if ok {
			choiceSelected = "X"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, choiceSelected, choice)
	}

	s += "\nPress q to quit.\n"

	return tea.NewView(s)
}