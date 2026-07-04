package room

import (
	"strings"
	"sync"
	"time"
)

type MemoryStore struct {
	mu    sync.Mutex
	rooms map[RoomCode]Room
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms: make(map[RoomCode]Room),
	}
}

func (m *MemoryStore) AuthenticateSession(code RoomCode, token SessionToken) (SessionContext, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exists := m.rooms[code]
	if !exists {
		return SessionContext{}, ErrInactiveRoom
	}

	if time.Now().After(room.ExpiresAt) {
		return SessionContext{}, ErrInactiveRoom
	}

	session, err := retrieveSession(token, room.Sessions)
	if err != nil {
		return SessionContext{}, err
	}

	participant, err := retrieveParticipant(session.ParticipantID, room.Participants)
	if err != nil {
		return SessionContext{}, err
	}

	return SessionContext{
		Code:          room.Code,
		ParticipantID: participant.ID,
		DisplayName:   participant.DisplayName,
		ExpiresAt:     room.ExpiresAt,
	}, nil
}

func (m *MemoryStore) CreateRoom(room Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rooms[room.Code]; exists {
		return ErrRoomCodeAlreadyExists
	}

	m.rooms[room.Code] = room
	return nil
}

func (m *MemoryStore) JoinRoom(code RoomCode, participant Participant, session Session) (Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, exists := m.rooms[code]
	if !exists {
		return Room{}, ErrInactiveRoom
	}

	if time.Now().After(room.ExpiresAt) {
		return Room{}, ErrInactiveRoom
	}

	if len(room.Participants) >= 8 {
		return Room{}, ErrRoomFull
	}

	if isNameAlreadyTaken(participant.DisplayName, room.Participants) {
		return Room{}, ErrDisplayNameTaken
	}

	room.Participants = append(room.Participants, participant)
	room.Sessions = append(room.Sessions, session)

	m.rooms[code] = room

	return m.rooms[code], nil
}

func isNameAlreadyTaken(displayName string, participants []Participant) bool {
	name := strings.ToLower(strings.TrimSpace(displayName))

	for _, taken := range participants {
		if strings.ToLower(strings.TrimSpace(taken.DisplayName)) == name {
			return true
		}
	}

	return false
}

func retrieveSession(st SessionToken, sessions []Session) (Session, error) {
	for _, session := range sessions {
		if st == session.Token {
			return session, nil
		}
	}

	return Session{}, ErrInvalidSession
}

func retrieveParticipant(participantID ParticipantID, participants []Participant) (Participant, error) {
	for _, participant := range participants {
		if participantID == participant.ID {
			return participant, nil
		}
	}

	return Participant{}, ErrInvalidSession
}
