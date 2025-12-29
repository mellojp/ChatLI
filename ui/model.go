package ui

import (
	"chatli/api"
	"chatli/data"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type uiState int

const (
	loginView uiState = iota
	roomListView
	chatView
	joinRoomView
)

type Model struct {
	State        uiState
	ChatsHistory map[string][]data.Message
	TextArea     textarea.Model
	Cursor       int
	CurrentRoom  string
	Session      data.Session
	WindowHeight int
	WindowWidth  int
	ErrorMsg     string
	SocketsConns map[string]*websocket.Conn
}

func NewModel() *Model {
	ta := textarea.New()
	ta.Placeholder = "Digite um nome..."
	ta.Focus()
	ta.Prompt = ""
	ta.ShowLineNumbers = false

	ta.SetWidth(69)
	ta.SetHeight(1)

	ta.CharLimit = 69

	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.BlurredStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	return &Model{
		State:        loginView,
		ChatsHistory: map[string][]data.Message{},
		TextArea:     ta,
		Cursor:       0,
		CurrentRoom:  "",
		ErrorMsg:     "",
		SocketsConns: map[string]*websocket.Conn{},
	}
}

func (*Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.State == loginView {
				name := m.TextArea.Value()
				session, err := api.CreateSession(name)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.Session = *session
				m.ErrorMsg = ""
				m.State = roomListView
				m.TextArea.Reset()
				return m, cmd
			}
			if m.State == roomListView {
				if len(m.Session.JoinedRooms) == 0 {
					return m, nil
				}

				m.CurrentRoom = m.Session.JoinedRooms[m.Cursor]
				m.State = chatView
				m.TextArea.Reset()

				conn, ok := m.SocketsConns[m.CurrentRoom]
				if !ok || conn == nil {
					wsConn, err := api.ConnectWebSocket(m.Session, m.CurrentRoom)
					if err != nil {
						m.ErrorMsg = err.Error()
						return m, nil
					}
					m.SocketsConns[m.CurrentRoom] = wsConn
					conn = wsConn
				}

				v, err := api.LoadChatMessages(m.Session, m.CurrentRoom, 50)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}
				m.ChatsHistory[m.CurrentRoom] = v

				return m, WaitForMessage(conn)
			}
			if m.State == joinRoomView {
				roomId := m.TextArea.Value()
				err := api.JoinRoom(m.Session, roomId)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				wsConn, err := api.ConnectWebSocket(m.Session, roomId)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.SocketsConns[roomId] = wsConn
				cmd := WaitForMessage(wsConn)
				v, err := api.LoadChatMessages(m.Session, roomId, 50)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.ChatsHistory[roomId] = v

				m.ErrorMsg = ""
				m.Session.JoinedRooms = append(m.Session.JoinedRooms, roomId)
				m.State = chatView
				m.CurrentRoom = roomId
				m.TextArea.Reset()
				return m, cmd
			}
			if m.State == chatView {
				m.TextArea.Placeholder = "Digite sua mensagem..."
				conn := m.SocketsConns[m.CurrentRoom]
				msg := data.Message{
					//Id: ,gerar uuid aleatório (string)
					Type:      "chat",
					User:      m.Session.Username,
					Timestamp: time.Now(),
					Content:   m.TextArea.Value(),
					RoomId:    m.CurrentRoom,
				}
				err := conn.WriteJSON(msg)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.ChatsHistory[m.CurrentRoom] = append(m.ChatsHistory[m.CurrentRoom], msg)
				m.TextArea.Reset()
				return m, cmd
			}
		case "n":
			if m.State == roomListView {
				newRoom, err := api.CreateRoom(m.Session)
				roomId := (*newRoom).Id
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				wsConn, err := api.ConnectWebSocket(m.Session, roomId)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.SocketsConns[roomId] = wsConn
				cmd := WaitForMessage(wsConn)
				v, err := api.LoadChatMessages(m.Session, roomId, 50)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, cmd
				}
				m.ChatsHistory[roomId] = v

				m.ErrorMsg = ""
				m.Session.JoinedRooms = append(m.Session.JoinedRooms, roomId)
				m.TextArea.Reset()
				return m, cmd
			}
		case "e":
			if m.State == roomListView {
				m.State = joinRoomView
				m.TextArea.Reset()
				return m, cmd
			}
		case "up", "k":
			if m.Cursor > 0 && m.State == roomListView {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Session.JoinedRooms)-1 && m.State == roomListView {
				m.Cursor++
			}
		case "esc":
			m.ErrorMsg = ""
			switch m.State {
			case loginView:
				m.State = loginView
			case roomListView:
				m.State = loginView
			case joinRoomView:
				m.State = roomListView
			case chatView:
				m.State = roomListView
			}
		}
	case tea.WindowSizeMsg:
		m.WindowHeight = msg.Height
		m.WindowWidth = msg.Width
		m.TextArea.SetWidth(70)

	case data.Message:
		m.ChatsHistory[msg.RoomId] = append(m.ChatsHistory[msg.RoomId], msg)
		cmd := WaitForMessage(m.SocketsConns[msg.RoomId])
		return m, cmd
	}
	if m.State != roomListView {
		m.TextArea, cmd = m.TextArea.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	s := ""
	switch m.State {
	case loginView:
		s = RenderLogin(m)
	case roomListView:
		s = RenderRoomsList(m)
	case joinRoomView:
		s = RenderJoinRoom(m)
	case chatView:
		s = RenderChatView(m)
	}
	temp := AppStyle.Render(s)
	r := lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, temp)
	return r
}

func WaitForMessage(conn *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		msg := data.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			return err
		}
		return msg
	}
}
