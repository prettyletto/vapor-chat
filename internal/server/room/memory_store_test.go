package room

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestMemoryStoreCreateRoomSucceeds(t *testing.T) {
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

func TestMemoryStoreJoinRoomSucceeds(t *testing.T) {
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

	participant := Participant{
		ID:          ParticipantID("participant-1"),
		DisplayName: "Joe123",
	}
	session := Session{
		Token:         SessionToken("token-2"),
		ParticipantID: ParticipantID("participant-1"),
	}

	updatedRoom, err := store.JoinRoom(room.Code, participant, session)
	if err != nil {
		t.Fatalf("expected first room join to succeed, got %v", err)
	}

	if len(updatedRoom.Participants) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(updatedRoom.Participants))
	}

	if len(updatedRoom.Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(updatedRoom.Sessions))
	}

	if updatedRoom.Participants[1].DisplayName != "Joe123" {
		t.Fatalf("expected joined participant display name %q, got %q", "Joe123", updatedRoom.Participants[1].DisplayName)
	}

	if updatedRoom.Sessions[1].Token != SessionToken("token-2") {
		t.Fatalf("expected joined session token %q, got %q", SessionToken("token-2"), updatedRoom.Sessions[1].Token)
	}
}

func TestMemoryStoreJoinRoomRejectsDuplicateName(t *testing.T) {
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

	participant := Participant{
		ID:          ParticipantID("participant-1"),
		DisplayName: "alex123",
	}
	session := Session{
		Token:         SessionToken("token-2"),
		ParticipantID: ParticipantID("participant-1"),
	}

	_, err = store.JoinRoom(room.Code, participant, session)
	if !errors.Is(err, ErrDisplayNameTaken) {
		t.Fatalf("expected ErrDisplayNameTaken, got %v", err)
	}
}

func TestMemoryStoreJoinRoomRejectsExpiredRoom(t *testing.T) {
	store := NewMemoryStore()

	room := Room{
		Code:                 RoomCode("ABCDEFGHIJ"),
		ExpiresAt:            time.Now().Add(-30 * time.Minute),
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

	participant := Participant{
		ID:          ParticipantID("participant-1"),
		DisplayName: "joe123",
	}
	session := Session{
		Token:         SessionToken("token-2"),
		ParticipantID: ParticipantID("participant-1"),
	}

	_, err = store.JoinRoom(room.Code, participant, session)
	if !errors.Is(err, ErrInactiveRoom) {
		t.Fatalf("expected ErrInactiveRoom, got %v", err)
	}
}

func TestMemoryStoreJoinRoomRejectsFullRoom(t *testing.T) {
	store := NewMemoryStore()
	participants := make([]Participant, 0, 8)
	sessions := make([]Session, 0, 8)

	for i := range 8 {
		participants = append(participants, Participant{
			ID:          ParticipantID(fmt.Sprintf("create-%d", i)),
			DisplayName: fmt.Sprintf("Alex%d", i),
		})
		sessions = append(sessions, Session{
			Token:         SessionToken(fmt.Sprintf("session-%d", i)),
			ParticipantID: ParticipantID(fmt.Sprintf("create-%d", i)),
		})
	}

	room := Room{
		Code:                 RoomCode("ABCDEFGHIJ"),
		ExpiresAt:            time.Now().Add(30 * time.Minute),
		CreatorParticipantID: ParticipantID("creator-1"),
		Participants:         participants,
		Sessions:             sessions,
	}

	err := store.CreateRoom(room)
	if err != nil {
		t.Fatalf("expected first room creation to succeed, got %v", err)
	}

	participant := Participant{
		ID:          ParticipantID("participant-1"),
		DisplayName: "joe123",
	}
	session := Session{
		Token:         SessionToken("token-2"),
		ParticipantID: ParticipantID("participant-1"),
	}

	_, err = store.JoinRoom(room.Code, participant, session)
	if !errors.Is(err, ErrRoomFull) {
		t.Fatalf("expected ErrRoomFull, got %v", err)
	}
}
