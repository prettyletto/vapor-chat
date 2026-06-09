package room

import (
	"strings"
	"testing"
)

type fakeStore struct {
	createCalls int
	failUntil   int
}

func (f *fakeStore) CreateRoom(room Room) error {
	f.createCalls++
	if f.createCalls <= f.failUntil {
		return ErrRoomCodeAlreadyExists
	}
	return nil
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

	if result.Code == "" {
		t.Fatal("expected room code to be set")
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
		wantErrPart string
	}{
		{
			name: "empty after trim",
			input: CreateRoomInput{
				DisplayName: "   ",
				TTLPreset:   TTL30Minutes,
			},
			wantErrPart: "display name",
		},
		{
			name: "too short",
			input: CreateRoomInput{
				DisplayName: "ab",
				TTLPreset:   TTL30Minutes,
			},
			wantErrPart: "display name",
		},
		{
			name: "too long",
			input: CreateRoomInput{
				DisplayName: "abcdefghijklmnopqrstuvwxyz1234567",
				TTLPreset:   TTL30Minutes,
			},
			wantErrPart: "display name",
		},
		{
			name: "non alphanumeric",
			input: CreateRoomInput{
				DisplayName: "alex!",
				TTLPreset:   TTL30Minutes,
			},
			wantErrPart: "alphanumeric",
		},
		{
			name: "invalid ttl",
			input: CreateRoomInput{
				DisplayName: "Alex123",
				TTLPreset:   TTLPreset("9 hours"),
			},
			wantErrPart: "ttl",
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

			if !strings.Contains(strings.ToLower(err.Error()), tt.wantErrPart) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrPart, err.Error())
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

	if result.Code == "" {
		t.Fatal("expected room code to be set")
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
