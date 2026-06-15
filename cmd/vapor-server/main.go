package main

import (
	"log"
	"net/http"

	serverhttp "github.com/prettyletto/vapor-chat/internal/server/http"
	"github.com/prettyletto/vapor-chat/internal/server/room"
)

func main() {
	store := room.NewMemoryStore()
	service := room.NewService(store)
	roomHandler := serverhttp.NewRoomHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /rooms", roomHandler.CreateRoom)
	mux.HandleFunc("POST /rooms/join", roomHandler.JoinRoom)

	addr := ":8080"

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
