package components

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type RoomData struct {
	Id         string
	UsersCount int
}

func (r RoomData) FilterValue() string { return r.Id }

// NOTE: On WindowSizeMsg should you might consider calling SetWidth, SetHeight
type RoomListModel struct {
	list list.Model
}

func NewRoomListModel(rooms []RoomData, width int, height int) RoomListModel {
	room_items := make([]list.Item, 0, len(rooms))
	for _, room := range rooms {
		room_items = append(room_items, room)
	}
	li := list.New(room_items, itemDelegate{}, width, height)
	li.Title = "All Open Rooms"

	return RoomListModel{
		list: li,
	}
}

func (m RoomListModel) Init() tea.Cmd {
	return nil
}

func (m RoomListModel) Update(msg tea.Msg) (RoomListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RoomListModel) View() string {
	return m.list.View()
}

func (m *RoomListModel) SetHeight(h int) {
	m.list.SetHeight(h)
}

func (m *RoomListModel) SetWidth(w int) {
	m.list.SetWidth(w)
}

type itemDelegate struct{}

func (i itemDelegate) Render(w io.Writer, li list.Model, index int, listItem list.Item) {
	item, ok := listItem.(RoomData)
	if !ok {
		return
	}

	roomName := item.Id
	if index == li.Index() {
		roomName += "-> "
	}
	fmt.Fprintf(w, "%s %d", roomName, item.UsersCount)
}

func (i itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (i itemDelegate) Height() int                               { return 1 }
func (i itemDelegate) Spacing() int                              { return 2 }
