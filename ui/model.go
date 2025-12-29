package ui

import (
	"chatli/api"
	"chatli/data"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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

type SocketError struct {
	RoomId string
	Err    error
}

type ReconnectMsg struct {
	RoomId string
}

type Model struct {
	State        uiState
	ChatsHistory map[string][]data.Message
	TextArea     textarea.Model
	Viewport     viewport.Model
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

	vp := viewport.New(70, 20)
	vp.SetContent("Carregando mensagens...")

	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.BlurredStyle.Base = lipgloss.NewStyle()
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	return &Model{
		State:        loginView,
		ChatsHistory: map[string][]data.Message{},
		TextArea:     ta,
		Viewport:     vp,
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
	var cmd, vpCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.State {
			case loginView:
				name := m.TextArea.Value()
				session, err := api.CreateSession(name)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}
				m.Session = *session
				m.ErrorMsg = ""
				m.State = roomListView
				m.TextArea.Reset()
				return m, nil

			case roomListView:
				if len(m.Session.JoinedRooms) == 0 {
					return m, nil
				}

				m.CurrentRoom = m.Session.JoinedRooms[m.Cursor]
				m.State = chatView
				m.TextArea.Reset()
				err := api.JoinRoom(m.Session, m.CurrentRoom)
				if err != nil {
					m.ErrorMsg = "Erro ao entrar na sala: " + err.Error()
					return m, nil
				}
				v, err := api.LoadChatMessages(m.Session, m.CurrentRoom, 50)
				if err != nil {
					m.ErrorMsg = err.Error()
				} else {
					m.ChatsHistory[m.CurrentRoom] = v
					m.Viewport.SetContent(RenderChatView(m))
					m.Viewport.GotoBottom()
				}
				conn, ok := m.SocketsConns[m.CurrentRoom]
				if !ok || conn == nil {
					wsConn, err := api.ConnectWebSocket(m.Session, m.CurrentRoom)
					if err != nil {
						m.ErrorMsg = err.Error()
						return m, nil
					}
					m.SocketsConns[m.CurrentRoom] = wsConn
					return m, WaitForMessage(wsConn, m.CurrentRoom)
				}
				return m, nil

			case joinRoomView:
				roomId := m.TextArea.Value()
				err := api.JoinRoom(m.Session, roomId)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}
				wsConn, err := api.ConnectWebSocket(m.Session, roomId)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}
				m.SocketsConns[roomId] = wsConn

				v, err := api.LoadChatMessages(m.Session, roomId, 50)
				if err != nil {
					m.ErrorMsg = err.Error()
				} else {
					m.ChatsHistory[roomId] = v
					m.Viewport.SetContent(RenderChatView(m))
					m.Viewport.GotoBottom()
				}

				m.ErrorMsg = ""
				m.Session.JoinedRooms = append(m.Session.JoinedRooms, roomId)
				m.State = chatView
				m.CurrentRoom = roomId
				m.TextArea.Reset()
				return m, WaitForMessage(wsConn, roomId)

			case chatView:
				content := m.TextArea.Value()
				if content == "" {
					return m, nil
				}

				conn := m.SocketsConns[m.CurrentRoom]
				msg := data.Message{
					Type:      "chat",
					User:      m.Session.Username,
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
					Content:   content,
					RoomId:    m.CurrentRoom,
				}

				err := conn.WriteJSON(msg)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}
				m.TextArea.Reset()
				return m, nil
			}
		case "n":
			if m.State == roomListView {
				newRoom, err := api.CreateRoom(m.Session)
				if err != nil {
					m.ErrorMsg = err.Error()
					return m, nil
				}

				// Apenas adiciona à lista e limpa mensagens de erro
				m.Session.JoinedRooms = append(m.Session.JoinedRooms, newRoom.Id)
				m.ErrorMsg = ""
				return m, nil
			}
		case "e":
			if m.State == roomListView {
				m.State = joinRoomView
				m.TextArea.Reset()
				return m, cmd
			}
		case "up":
			if m.Cursor > 0 && m.State == roomListView {
				m.Cursor--
			}
		case "down":
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
		m.Viewport.Height = msg.Height - 14
		m.Viewport.Width = msg.Width

	case data.Message:
		m.ChatsHistory[msg.RoomId] = append(m.ChatsHistory[msg.RoomId], msg)
		s := RenderChatView(m)
		m.Viewport.SetContent(s)
		m.Viewport.GotoBottom()
		cmd := WaitForMessage(m.SocketsConns[msg.RoomId], msg.RoomId)
		return m, cmd

	case SocketError:
		m.ErrorMsg = "Erro detetado: " + msg.Err.Error()
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return ReconnectMsg{RoomId: msg.RoomId}
		})
	case ReconnectMsg:
		newConn, err := api.ConnectWebSocket(m.Session, msg.RoomId)
		if err != nil {
			m.ErrorMsg = "Erro na reconexão: " + err.Error()
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return ReconnectMsg{RoomId: msg.RoomId}
			})
		}
		m.SocketsConns[msg.RoomId] = newConn
		m.ErrorMsg = ""
		return m, WaitForMessage(newConn, msg.RoomId)
	}
	if m.State != roomListView {
		m.TextArea, cmd = m.TextArea.Update(msg)
		m.Viewport, vpCmd = m.Viewport.Update(msg)
	}

	return m, tea.Batch(cmd, vpCmd)
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
		header := RoomTitleStyle.Render("Sala: " + m.CurrentRoom)
		body := m.Viewport.View()
		footer := InputStyle.Render(m.TextArea.View())
		if m.ErrorMsg != "" {
			errMsg := ErrorStyle.Render(m.ErrorMsg)
			s = lipgloss.JoinVertical(lipgloss.Left, header, body, errMsg, footer)
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
		}
	}
	temp := AppStyle.Render(s)
	r := lipgloss.Place(m.WindowWidth, m.WindowHeight, lipgloss.Center, lipgloss.Center, temp)
	return r
}

func WaitForMessage(conn *websocket.Conn, roomId string) tea.Cmd {
	return func() tea.Msg {
		msg := data.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			return SocketError{RoomId: roomId, Err: err}
		}
		return msg
	}
}
