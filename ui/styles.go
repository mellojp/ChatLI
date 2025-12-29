package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var ColorMap = map[int]string{
	1: "#d42222ff",
	2: "#1411c4ff",
	3: "#0aaf05ff",
	4: "#c6d317ff",
	5: "#ca1ddaff",
	6: "#50f5dfff",
}

var AppStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	BorderForeground(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Border(lipgloss.RoundedBorder()).
	Margin(1, 1).
	Width(80)

var InputStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("#ec60ecff")).
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	MarginTop(1)

var SenderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5")).
	Bold(true)

var TimeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	Faint(true)

var RoomTitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7D56F4")).
	Bold(true).
	Border(lipgloss.NormalBorder(), false, false, true, false).
	Padding(0, 1).
	MarginBottom(1)

var UnselectedItemStyle = lipgloss.NewStyle().
	PaddingLeft(2)

var SelectedItemStyle = lipgloss.NewStyle().
	PaddingLeft(4).
	Foreground(lipgloss.Color("#ec60ecff")).
	Bold(true)

var ButtonStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#555")).
	Padding(0, 2).
	MarginRight(1)

var ActiveButtonStyle = ButtonStyle.Copy().
	BorderForeground(lipgloss.Color("#ec60ecff")).
	Foreground(lipgloss.Color("#ec60ecff")).
	Bold(true)

var ListHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7D56F4")).
	Bold(true).
	MarginTop(1).
	MarginBottom(1)

var HelpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	MarginTop(2).
	MarginBottom(2)

var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196")).
	Bold(true).
	PaddingTop(1)

func HashColor(username string, colors map[int]string) string {
	var hashCode int
	for _, v := range username {
		hashCode += int(v)
	}
	hashCode = (hashCode+37)%len(colors) + 1
	return colors[hashCode]
}
