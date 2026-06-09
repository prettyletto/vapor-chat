package room

import (
	"errors"
	"testing"
	"time"
)

func TestMemoryStoreCreateRoomSuceeds(t *testing.T) {
	store := NewMemoryStore()

	room := Room{
		Code:                 RoomCode("ABCDEFGHIJ"),
		ExpiresAt:            time.Now().Add(30 * time.Minute),
		CreatorParticipantID: ParticipantID("creator-1"),
		Participants: []Participant{
			{
				ID:          ParticipantID("creator-1"),
				DisplayName: "Alex123",
			},
		}, Sessions: []Session{
			{
				Token:         SessionToken("token-1"),
				ParticipantID: ParticipantID("creator-1"),
			},
		},
	}

	err := store.CreateRoom(room)
	if err != nil {
		t.Fatalf("expcted room creation to succeed, got %v", err)
	}
}

func TestMemoryStoreCreateRoomRejectsDuplicateCode(t *testing.T) {
	store := NewMemoryStore()

	room := Room{
		Code:                 RoomCode("ABCDEFGHIJ"),
		ExpiresAt:            time.Now().Add(30 * time.Minute),
		CreatorParticipantID: ParticipantID("creator-1"),
		Participants: []Participant{
			{
				ID:          ParticipantID("creator-1"),
				DisplayName: "Alex123",
			},
		},
		Sessions: []Session{
			{
				Token:         SessionToken("token-1"),
				ParticipantID: ParticipantID("creator-1"),
			},
		},
	}

	err := store.CreateRoom(room)
	if err != nil {
		t.Fatalf("expected first room creation to succeed, got %v", err)
	}

	err = store.CreateRoom(room)
	if !errors.Is(err, ErrRoomCodeAlreadyExists) {
		t.Fatalf("expected ErrRoomCodeAlreadyExists, got %v", err)
	}
}
