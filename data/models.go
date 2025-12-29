package data

import (
	"time"
)

type Message struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	User      string    `json:"user"`
	Content   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	RoomId    string    `json:"room_id"`
}

type Room struct {
	Id           string    `json:"room_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	ActiveUsers  []string  `json:"active_users"`
}

type Session struct {
	Id           string    `json:"session_id"`
	Username     string    `json:"username"`
	JoinedRooms  []string  `json:"joined_rooms"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}
