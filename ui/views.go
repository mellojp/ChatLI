package ui

import (
	"fmt"
	"strings"

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
	centeredAscii := lipgloss.PlaceHorizontal(m.WindowWidth, lipgloss.Center, SystemStyle.Render(ascii))
	s := centeredAscii + "\n\n"
	s += SystemStyle.Render("user@terminal:~$ login") + "\n"
	s += SystemStyle.Render("Username: ") + InputStyle.Render(m.TextArea.View())
	if m.ErrorMsg != "" {
		s += "\n" + ErrorStyle.Render("error: "+m.ErrorMsg)
	}
	return s
}

func RenderRoomsList(m *Model) string {
	s := fmt.Sprintf("%s@terminal:~$ /chatli/rooms ls", m.Session.Username)
	s = SystemStyle.Render(s) + "\n"
	for i, v := range m.Session.JoinedRooms {
		if i == m.Cursor {
			s += SelectedItemStyle.Render("> "+v) + "\n"
		} else {
			s += UnselectedItemStyle.Render("  "+v) + "\n"
		}
	}
	s += "\n" + HelpStyle.Render("[n] new | [e] join | [enter] select | [esc] logout")
	if m.ErrorMsg != "" {
		s += "\n" + ErrorStyle.Render("error: "+m.ErrorMsg)
	}
	return s
}

func RenderJoinRoom(m *Model) string {
	s := fmt.Sprintf("%s@terminal:~$ join-room --id", m.Session.Username)
	s = SystemStyle.Render(s) + "\n"
	s += SystemStyle.Render("Room ID: ") + InputStyle.Render(m.TextArea.View())
	if m.ErrorMsg != "" {
		s += "\n" + ErrorStyle.Render("error: "+m.ErrorMsg)
	}
	return s
}

func RenderChatView(m *Model) string {
	m.TextArea.Placeholder = ""
	var b strings.Builder
	renderWidth := m.WindowWidth - 2

	for _, val := range m.ChatsHistory[m.CurrentRoom] {
		if val.Content == "" {
			continue
		}

		displayTime := val.Timestamp
		if len(displayTime) >= 16 {
			displayTime = displayTime[11:16]
		}

		//userColor := lipgloss.Color(HashColor(val.User, ColorMap))
		styledUser := SenderStyle.Render(val.User)

		if m.Session.Username == val.User {
			styledUser := SenderStyle.Render("você")
			line := fmt.Sprintf("%s <%s> [%s]", val.Content, styledUser, TimeStyle.Render(displayTime))
			b.WriteString(lipgloss.NewStyle().Width(renderWidth).Align(lipgloss.Right).Render(line) + "\n")
		} else {
			line := fmt.Sprintf("[%s] <%s> %s", TimeStyle.Render(displayTime), styledUser, val.Content)
			b.WriteString(lipgloss.NewStyle().Width(renderWidth).Align(lipgloss.Left).Render(line) + "\n")
		}
	}
	return b.String()
}
