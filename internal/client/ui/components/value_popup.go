package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PopUpModel struct {
	Value string

	ti textinput.Model
}

func NewValuePopUp(
	name string,
	maxLength int,
	validate textinput.ValidateFunc,
) PopUpModel {
	ti := textinput.New()
	ti.Placeholder = name
	ti.CharLimit = maxLength
	ti.Validate = validate
	ti.Width = maxLength
	ti.Focus()
	return PopUpModel{
		ti: ti,
	}
}

func (m PopUpModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PopUpModel) Update(msg tea.Msg) (PopUpModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.ti.Err != nil {
				break
			}
			m.Value = m.ti.Value()
			return m, nil
		}
	}

	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m PopUpModel) View() string {
	var err_msg string
	err := m.ti.Err
	if err != nil {
		err_msg = err.Error()
	}
	box := lipgloss.NewStyle().
		Height(7).
		Width(m.ti.CharLimit+8).
		Border(lipgloss.NormalBorder()).
		Align(0.5, 0.5)

	err_style := lipgloss.NewStyle().Foreground(lipgloss.Color("31"))

	return box.Render(
		fmt.Sprintf(
			"%s\n\n%s",
			m.ti.View(),
			err_style.Render(err_msg),
		),
	)
}
