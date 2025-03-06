package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/starathel/gchat/internal/client/ui/components"
)

func StartBubbleTea() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type appState int

const (
	InputLoginState appState = iota
	RoomsListState
)

type model struct {
	username string

	state         appState
	usernamePopUp components.PopUpModel
}

func newModel() model {
	return model{
		usernamePopUp: components.NewValuePopUp("Username", 20, nil),
	}
}

func (m model) Init() tea.Cmd {
	return m.usernamePopUp.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	switch m.state {
	case InputLoginState:
		m.usernamePopUp, cmd = m.usernamePopUp.Update(msg)
		if m.usernamePopUp.Value != "" {
			m.username = m.usernamePopUp.Value
			m.state = RoomsListState
		}
	case RoomsListState:
		break
	default:
		panic(fmt.Sprintf("Invalid state: %v", m.state))
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case InputLoginState:
		return m.usernamePopUp.View()
	case RoomsListState:
		return fmt.Sprintf("You are: %s", m.username)
	}
	panic(fmt.Sprintf("Invalid state: %v", m.state))
}
