package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const ascii = `
	 ________  ___  ___  ________  _________  ___       ___     
    |\   ____\|\  \|\  \|\   __  \|\___   ___\\  \     |\  \    
    \ \  \___|\ \  \\\  \ \  \|\  \|___ \  \_\ \  \    \ \  \   
     \ \  \    \ \   __  \ \   __  \   \ \  \ \ \  \    \ \  \  
      \ \  \____\ \  \ \  \ \  \ \  \   \ \  \ \ \  \____\ \  \ 
       \ \_______\ \__\ \__\ \__\ \__\   \ \__\ \ \_______\ \__\
        \|_______|\|__|\|__|\|__|\|__|    \|__|  \|_______|\|__|
`

func RenderLogin(m *Model) string {
	s := ascii + "\n"
	s += InputStyle.Render(m.TextArea.View())
	s += ErrorStyle.Render(m.ErrorMsg)
	return s
}

func RenderRoomsList(m *Model) string {
	s := "Salas disponívies:\n\n"
	for i, v := range m.Session.JoinedRooms {
		if i == m.Cursor {
			line := "> " + v
			s += SelectedItemStyle.Render(line) + "\n\n"
		} else {
			line := "  " + v
			s += UnselectedItemStyle.Render(line) + "\n\n"
		}
	}
	s += HelpStyle.Render("Nova Sala [n] • Entrar em Sala [e] • Sair [esc]")
	s += ErrorStyle.Render(m.ErrorMsg)
	return s
}

func RenderJoinRoom(m *Model) string {
	s := "Digite o código da sala\n"
	s += InputStyle.Render(m.TextArea.View())
	s += ErrorStyle.Render(m.ErrorMsg)
	return s
}

func RenderChatView(m *Model) string {
	s := "Sala: " + m.CurrentRoom + "\n\n"
	for _, val := range m.ChatsHistory[m.CurrentRoom] {
		if val.Content != "" {
			time := val.Timestamp.Format("15:04")
			styledTime := TimeStyle.Render(time)
			styledUser := SenderStyle.Render(val.User)
			if m.Session.Username == val.User {
				line := fmt.Sprintf("%s [%s]\n %s", styledUser, styledTime, val.Content)
				s += lipgloss.NewStyle().Width(76).PaddingRight(3).Align(lipgloss.Right).Render(line) + "\n"
			} else {
				line := fmt.Sprintf("[%s] %s\n %s", styledTime, styledUser, val.Content)
				s += lipgloss.NewStyle().Width(76).PaddingRight(3).Align(lipgloss.Left).Render(line) + "\n"
			}
		}
	}
	s += ErrorStyle.Render(m.ErrorMsg)
	s += "\n" + InputStyle.Render(m.TextArea.View())
	return s
}
