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
	ConnectedToRoomState
)

type roomConnectionMsg int

const (
	ConnectedSuccessfully roomConnectionMsg = iota
)

type model struct {
	username string

	state appState
	w, h  int

	roomId string

	usernamePopUp components.PopUpModel
	roomsList     components.RoomListModel
}

func newModel() model {
	rooms := []components.RoomData{
		{Id: "Room 1", UsersCount: 12},
		{Id: "Room 2", UsersCount: 13},
		{Id: "Aboba", UsersCount: 69},
	}
	return model{
		usernamePopUp: components.NewValuePopUp("Username", 20, nil),
		roomsList:     components.NewRoomListModel(rooms, 10, 10),
	}
}

func (m model) Init() tea.Cmd {
	return m.usernamePopUp.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		// TODO
		m.roomsList.SetHeight(m.h - 20)
		m.roomsList.SetWidth(m.w - 20)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case roomConnectionMsg:
		switch msg {
		case ConnectedSuccessfully:
			m.state = ConnectedToRoomState
			return m, nil
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
		m.roomsList, cmd = m.roomsList.Update(msg)
		if m.roomsList.SelectedRoom != "" {
			m.roomId = m.roomsList.SelectedRoom
			return m, connectToRoom(m.roomId)
		}
	case ConnectedToRoomState:
		return m, cmd
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
		return m.roomsList.View()
	case ConnectedToRoomState:
		return fmt.Sprintf("%s connected to room %s", m.username, m.roomId)
	}
	panic(fmt.Sprintf("Invalid state: %v", m.state))
}

// TODO: all functions conserning connection to room should be in separate file
func connectToRoom(_ string) tea.Cmd {
	return func() tea.Msg {
		return ConnectedSuccessfully
	}
}
