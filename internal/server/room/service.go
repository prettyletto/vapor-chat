package room

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var displayNamePattern = regexp.MustCompile(`^[A-Za-z0-9]+$`)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Service struct {
	store Store
}

type CreateRoomInput struct {
	DisplayName string
	TTLPreset   TTLPreset
}

type CreateRoomResult struct {
	Code         RoomCode
	TTLPreset    TTLPreset
	ExpiresAt    time.Time
	SessionToken SessionToken
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateRoom(input CreateRoomInput) (CreateRoomResult, error) {
	var out CreateRoomResult

	name := strings.TrimSpace(input.DisplayName)

	if len(name) < 3 || len(name) > 32 {
		return out, ValidationError{
			Field:   "display_name",
			Message: "must between 3 and 32 characters",
		}
	}

	if !displayNamePattern.MatchString(name) {
		return out, ValidationError{
			Field:   "display_name",
			Message: "must contain only alphanumerical characters",
		}
	}

	expiresAt, err := expirationFromPreset(input.TTLPreset)
	if err != nil {
		return out, err
	}

	for range 5 {
		rc, err := generateRoomCode()
		if err != nil {
			return out, err
		}
		creatorID := ParticipantID(uuid.NewString())
		sessionToken := SessionToken(uuid.NewString())

		participant := Participant{
			ID:          creatorID,
			DisplayName: name,
		}

		session := Session{
			Token:         sessionToken,
			ParticipantID: creatorID,
		}

		room := Room{
			Code:                 rc,
			ExpiresAt:            expiresAt,
			CreatorParticipantID: creatorID,
			Participants:         []Participant{participant},
			Sessions:             []Session{session},
		}
		outErr := s.store.CreateRoom(room)
		if outErr == nil {
			out.Code = rc
			out.ExpiresAt = expiresAt
			out.TTLPreset = input.TTLPreset
			out.SessionToken = sessionToken

			return out, nil
		}

		if outErr == ErrRoomCodeAlreadyExists {
			continue
		}

		return out, outErr
	}

	return out, fmt.Errorf("could not create room after 5 code collisions")
}

func generateRoomCode() (RoomCode, error) {
	code := make([]byte, 10)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		code[i] = alphabet[n.Int64()]
	}
	return RoomCode(code), nil
}

func expirationFromPreset(ttl TTLPreset) (time.Time, error) {
	switch ttl {
	case TTL15Minutes:
		return time.Now().Add(time.Minute * 15), nil
	case TTL30Minutes:
		return time.Now().Add(time.Minute * 30), nil
	case TTL1Hour:
		return time.Now().Add(time.Minute * 60), nil
	case TTL2Hours:
		return time.Now().Add(time.Minute * 120), nil
	default:
		return time.Time{}, ValidationError{
			Field:   "ttl_preset",
			Message: "is not valid",
		}
	}
}
