package protocol

import "time"

type EventType string

const (
	EventCreateRoom  EventType = "create_room"
	EventJoinRoom    EventType = "join_room"
	EventChatMessage EventType = "chat_message"

	EventMemberJoined EventType = "member_joined"
	EventMemberLeft   EventType = "member_left"

	EventRoomCreated EventType = "room_created"
	EventRoomExpired EventType = "room_expired"
	EventRoomJoined  EventType = "room_joined"

	Error EventType = "error"
)

type Event struct {
	Type      EventType `json:"type"`
	RoomID    string    `json:"room_id,omitempty"`
	SenderID  string    `json:"sender_id,omitempty"`
	MemberID  string    `json:"member_id,omitempty"`
	Text      string    `json:"text,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
