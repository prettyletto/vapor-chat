package room

import "time"

type (
	RoomCode      string
	SessionToken  string
	TTLPreset     string
	ParticipantID string
)

const (
	TTL15Minutes TTLPreset = "15 minutes"
	TTL30Minutes TTLPreset = "30 minutes"
	TTL1Hour     TTLPreset = "1 hour"
	TTL2Hours    TTLPreset = "2 hours"
)

type Participant struct {
	ID          ParticipantID
	DisplayName string
}

type Session struct {
	Token         SessionToken
	ParticipantID ParticipantID
}

type Room struct {
	Code                 RoomCode
	ExpiresAt            time.Time
	CreatorParticipantID ParticipantID
	Participants         []Participant
	Sessions             []Session
}

type SessionContext struct {
	Code          RoomCode
	ParticipantID ParticipantID
	DisplayName  string
	ExpiresAt     time.Time
}
