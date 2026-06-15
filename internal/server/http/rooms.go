package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	serverRoom "github.com/prettyletto/vapor-chat/internal/server/room"
)

type RoomHandler struct {
	service *serverRoom.Service
}

func NewRoomHandler(service *serverRoom.Service) *RoomHandler {
	return &RoomHandler{service: service}
}

type createRoomRequest struct {
	DisplayName string `json:"display_name"`
	TTLPreset   string `json:"ttl_preset"`
}

type createRoomResponse struct {
	Code         string `json:"code"`
	TTLPreset    string `json:"ttl_preset"`
	ExpiresAt    string `json:"expires_at"`
	SessionToken string `json:"session_token"`
}

type joinRoomRequest struct {
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
}

type joinRoomResponse struct {
	Code         string `json:"code"`
	ExpiresAt    string `json:"expires_at"`
	SessionToken string `json:"session_token"`
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createRoomRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	result, err := h.service.CreateRoom(serverRoom.CreateRoomInput{
		DisplayName: req.DisplayName,
		TTLPreset:   serverRoom.TTLPreset(req.TTLPreset),
	})

	var validationErr serverRoom.ValidationError
	if err != nil {
		if errors.As(err, &validationErr) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, createRoomResponse{
		Code:         string(result.Code),
		TTLPreset:    string(result.TTLPreset),
		ExpiresAt:    result.ExpiresAt.Format(time.RFC3339),
		SessionToken: string(result.SessionToken),
	})
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req joinRoomRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	result, err := h.service.JoinRoom(serverRoom.JoinRoomInput{
		Code:        serverRoom.RoomCode(req.Code),
		DisplayName: req.DisplayName,
	})

	var validationErr serverRoom.ValidationError
	if err != nil {
		if errors.As(err, &validationErr) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, serverRoom.ErrInactiveRoom) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, serverRoom.ErrRoomFull) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, serverRoom.ErrDisplayNameTaken) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, joinRoomResponse{
		Code:         string(result.Code),
		ExpiresAt:    result.ExpiresAt.Format(time.RFC3339),
		SessionToken: string(result.SessionToken),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
