package socket

import (
	"fmt"
	"hangman/middlewares"
	"hangman/socket/events"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func Socket() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "create-room", func(s socketio.Conn, msg events.Room) {
		events.CreateRooms(msg)
	})

	server.OnEvent("/", "check-rooms", func(s socketio.Conn, roomName string) {
		s.Emit("do-room-exists", events.DoRoomAlreadyExists(roomName))
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, msg events.JoinRoom) {
		events.JoinRooms(msg)
		usersOfRoom := events.GetUserOfRooms(msg.RoomName)
		s.Emit("all-rooms", events.GetRooms())
		s.Emit("all-user-of-rooms", usersOfRoom)
	})

	http.Handle("/socket.io/", middlewares.CorsMiddleware(server))

	go server.Serve()
	defer server.Close()

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
