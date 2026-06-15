package room

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeStore struct {
	createCalls int
	failUntil   int
	joinRoom    Room
	joinErr     error
}

func (f *fakeStore) CreateRoom(room Room) error {
	f.createCalls++
	if f.createCalls <= f.failUntil {
		return ErrRoomCodeAlreadyExists
	}
	return nil
}

func (f *fakeStore) JoinRoom(code RoomCode, participant Participant, session Session) (Room, error) {
	if f.joinErr != nil {
		return Room{}, f.joinErr
	}
	return f.joinRoom, nil
}

func TestServiceCreateRoomSucceeds(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(store)

	input := CreateRoomInput{
		DisplayName: "Alex123",
		TTLPreset:   TTL30Minutes,
	}

	result, err := service.CreateRoom(input)
	if err != nil {
		t.Fatalf("expected CreateRoom to succeed, got %v", err)
	}

	if len(result.Code) != 6 {
		t.Fatalf("expected a 6-character room code, got %q", result.Code)
	}

	if result.SessionToken == "" {
		t.Fatal("expected session token to be set")
	}

	if result.TTLPreset != TTL30Minutes {
		t.Fatalf("expected ttl preset %q, got %q", TTL30Minutes, result.TTLPreset)
	}

	if result.ExpiresAt.IsZero() {
		t.Fatal("expected expiration timestamp to be set")
	}
}

func TestServiceCreateRoomValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateRoomInput
		wantField   string
		wantErrPart string
	}{
		{
			name: "empty after trim",
			input: CreateRoomInput{
				DisplayName: "   ",
				TTLPreset:   TTL30Minutes,
			},
			wantField:   ValidationFieldDisplayName,
			wantErrPart: "between 3 and 32 characters",
		},
		{
			name: "too short",
			input: CreateRoomInput{
				DisplayName: "ab",
				TTLPreset:   TTL30Minutes,
			},
			wantField:   ValidationFieldDisplayName,
			wantErrPart: "between 3 and 32 characters",
		},
		{
			name: "too long",
			input: CreateRoomInput{
				DisplayName: "abcdefghijklmnopqrstuvwxyz1234567",
				TTLPreset:   TTL30Minutes,
			},
			wantField:   ValidationFieldDisplayName,
			wantErrPart: "between 3 and 32 characters",
		},
		{
			name: "non alphanumeric",
			input: CreateRoomInput{
				DisplayName: "alex!",
				TTLPreset:   TTL30Minutes,
			},
			wantField:   ValidationFieldDisplayName,
			wantErrPart: "alphanumeric",
		},
		{
			name: "invalid ttl",
			input: CreateRoomInput{
				DisplayName: "Alex123",
				TTLPreset:   TTLPreset("9 hours"),
			},
			wantField:   ValidationFieldTTLPreset,
			wantErrPart: "not valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			service := NewService(store)

			_, err := service.CreateRoom(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var validationErr ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected ValidationError, got %T: %v", err, err)
			}

			if validationErr.Field != tt.wantField {
				t.Fatalf("expected validation field %q, got %q", tt.wantField, validationErr.Field)
			}

			if !strings.Contains(validationErr.Message, tt.wantErrPart) {
				t.Fatalf("expected validation message containing %q, got %q", tt.wantErrPart, validationErr.Message)
			}
		})
	}
}

func TestServiceCreateRoomRetriesOnCollision(t *testing.T) {
	store := &fakeStore{failUntil: 2}
	service := NewService(store)

	input := CreateRoomInput{
		DisplayName: "Alex123",
		TTLPreset:   TTL30Minutes,
	}

	result, err := service.CreateRoom(input)
	if err != nil {
		t.Fatalf("expected CreateRoom to succeed after retries, got %v", err)
	}

	if len(result.Code) != 6 {
		t.Fatalf("expected a 6-character room code, got %q", result.Code)
	}

	if store.createCalls != 3 {
		t.Fatalf("expected 3 create attempts, got %d", store.createCalls)
	}
}

func TestServiceCreateRoomFailsAfterFiveCollisions(t *testing.T) {
	store := &fakeStore{failUntil: 5}
	service := NewService(store)

	input := CreateRoomInput{
		DisplayName: "Alex123",
		TTLPreset:   TTL30Minutes,
	}

	_, err := service.CreateRoom(input)
	if err == nil {
		t.Fatal("expected error after five collisions, got nil")
	}
}

func TestServiceJoinRoomSucceeds(t *testing.T) {
	expiresAt := time.Now().Add(30 * time.Minute)

	store := &fakeStore{
		joinRoom: Room{
			Code:      RoomCode("ABC123"),
			ExpiresAt: expiresAt,
			Participants: []Participant{
				{ID: ParticipantID("creator-1"), DisplayName: "Alex123"},
				{ID: ParticipantID("participant-1"), DisplayName: "Joe123"},
			},
			Sessions: []Session{
				{Token: SessionToken("token-1"), ParticipantID: ParticipantID("creator-1")},
				{Token: SessionToken("token-2"), ParticipantID: ParticipantID("participant-1")},
			},
		},
	}
	service := NewService(store)

	input := JoinRoomInput{
		Code:        RoomCode("ABC123"),
		DisplayName: "Joe123",
	}

	result, err := service.JoinRoom(input)
	if err != nil {
		t.Fatalf("expected JoinRoom to succeed, got %v", err)
	}

	if result.Code != RoomCode("ABC123") {
		t.Fatalf("expected code %q, got %q", RoomCode("ABC123"), result.Code)
	}

	if result.ExpiresAt != expiresAt {
		t.Fatalf("expected expiresAt %v, got %v", expiresAt, result.ExpiresAt)
	}

	if result.SessionToken == "" {
		t.Fatal("expected session token to be set")
	}
}

func TestServiceJoinRoomValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       JoinRoomInput
		wantErrPart string
	}{
		{
			name: "empty after trim",
			input: JoinRoomInput{
				Code:        RoomCode("ABC123"),
				DisplayName: "   ",
			},
			wantErrPart: "display_name",
		},
		{
			name: "too short",
			input: JoinRoomInput{
				Code:        RoomCode("ABC123"),
				DisplayName: "ab",
			},
			wantErrPart: "display_name",
		},
		{
			name: "non alphanumeric",
			input: JoinRoomInput{
				Code:        RoomCode("ABC123"),
				DisplayName: "joe!",
			},
			wantErrPart: "alphanumeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&fakeStore{})

			_, err := service.JoinRoom(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(strings.ToLower(err.Error()), tt.wantErrPart) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrPart, err.Error())
			}
		})
	}
}

func TestServiceJoinRoomPropagatesInactiveRoom(t *testing.T) {
	store := &fakeStore{joinErr: ErrInactiveRoom}
	service := NewService(store)

	_, err := service.JoinRoom(JoinRoomInput{
		Code:        RoomCode("ABC123"),
		DisplayName: "Joe123",
	})
	if !errors.Is(err, ErrInactiveRoom) {
		t.Fatalf("expected ErrInactiveRoom, got %v", err)
	}
}

func TestServiceJoinRoomPropagatesRoomFull(t *testing.T) {
	store := &fakeStore{joinErr: ErrRoomFull}
	service := NewService(store)

	_, err := service.JoinRoom(JoinRoomInput{
		Code:        RoomCode("ABC123"),
		DisplayName: "Joe123",
	})
	if !errors.Is(err, ErrRoomFull) {
		t.Fatalf("expected ErrRoomFull, got %v", err)
	}
}

func TestServiceJoinRoomPropagatesDisplayNameTaken(t *testing.T) {
	store := &fakeStore{joinErr: ErrDisplayNameTaken}
	service := NewService(store)

	_, err := service.JoinRoom(JoinRoomInput{
		Code:        RoomCode("ABC123"),
		DisplayName: "Joe123",
	})
	if !errors.Is(err, ErrDisplayNameTaken) {
		t.Fatalf("expected ErrDisplayNameTaken, got %v", err)
	}
}
