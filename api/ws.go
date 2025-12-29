package api

import (
	"chatli/data"

	"github.com/gorilla/websocket"
)

const wsUrl = "wss://disposable-chat.onrender.com"

func ConnectWebSocket(s data.Session, roomId string) (*websocket.Conn, error) {
	connUrl := wsUrl + "/ws/" + roomId + "?session_id=" + s.Id
	socket, _, err := websocket.DefaultDialer.Dial(connUrl, nil)
	if err != nil {
		return nil, err
	}
	return socket, nil
}
