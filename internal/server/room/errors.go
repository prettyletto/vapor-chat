package room

import (
	"errors"
	"fmt"
)

const (
	ValidationFieldDisplayName = "display_name"
	ValidationFieldTTLPreset   = "ttl_preset"
)

var (
	ErrInactiveRoom     = errors.New("this room is inactive")
	ErrRoomFull         = errors.New("room full")
	ErrDisplayNameTaken = errors.New("display name already taken")
	ErrInvalidSession   = errors.New("invalid session")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s %s", e.Field, e.Message)
}
