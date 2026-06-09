package room

import "sync"

type MemoryStore struct {
	mu    sync.Mutex
	rooms map[RoomCode]Room
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms: make(map[RoomCode]Room),
	}
}

func (m *MemoryStore) CreateRoom( room Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _,exists := m.rooms[room.Code]; exists {
		return ErrRoomCodeAlreadyExists
	}

	m.rooms[room.Code] = room
	return nil
}
