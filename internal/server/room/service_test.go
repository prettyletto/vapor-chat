package room

import (
	"errors"
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
			wantErrPart: "alphanumerical",
		},
		{
			name: "invalid ttl",
			input: CreateRoomInput{
				DisplayName: "Alex123",
				TTLPreset:   TTLPreset("9 hours"),
			},
			wantField:   TTLPresetErr,
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
