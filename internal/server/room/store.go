package room

import "errors"

var ErrRoomCodeAlreadyExists = errors.New("Code Already exists for a active room")

type Store interface {
	CreateRoom(room Room) error
}
