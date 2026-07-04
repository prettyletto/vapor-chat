package room

import "errors"

var ErrRoomCodeAlreadyExists = errors.New("room code already exists")

type Store interface {
	CreateRoom(room Room) error
	JoinRoom(code RoomCode, participant Participant, session Session) (Room, error)
	AuthenticateSession(code RoomCode, token SessionToken) (SessionContext, error)
}
