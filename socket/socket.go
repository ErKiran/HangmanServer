package socket

import (
	"fmt"
	"hangman/middlewares"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type CreateGame struct {
	RoomName   string `json:"roomName"`
	SecretWord string `json:"secretWord"`
}

type JoinGame struct {
	NickName string `json:"nickName"`
	RoomName string `json:"roomName"`
}

type ConnectedUsers struct {
	RoomName string   `json:"roomName"`
	Players  []string `json:"players"`
}

func Socket() {

	var games []CreateGame
	var connectedUsers []ConnectedUsers

	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "create-game", func(s socketio.Conn, msg CreateGame) {
		alreadyExists := contains(games, msg.RoomName)
		if !alreadyExists {
			games = append(games, CreateGame{RoomName: msg.RoomName, SecretWord: msg.SecretWord})
		}
	})

	server.OnEvent("/", "check-rooms", func(s socketio.Conn, roomName string) {
		var isExists bool
		for _, name := range games {
			if name.RoomName == roomName {
				isExists = true
			}
		}
		s.Emit("do-room-exists", isExists)
	})

	server.OnEvent("/", "join-game", func(s socketio.Conn, msg JoinGame) {
		var players []string

		if len(connectedUsers) == 0 {
			connectedUsers = append(connectedUsers, ConnectedUsers{RoomName: msg.RoomName, Players: append(players, msg.NickName)})
		}
		for _, sockets := range connectedUsers {
			if sockets.RoomName == msg.RoomName {
				fmt.Println("Room Name Already Exists")
				sockets.Players = append(sockets.Players, msg.NickName)
				fmt.Println(sockets.Players)
			} else {
				connectedUsers = append(connectedUsers, ConnectedUsers{RoomName: msg.RoomName, Players: append(players, msg.NickName)})
			}
		}
		fmt.Println("ConnectedUsers", connectedUsers)
		s.Emit("all-games", games)
	})

	http.Handle("/socket.io/", middlewares.CorsMiddleware(server))

	go server.Serve()
	defer server.Close()

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func contains(s []CreateGame, e string) bool {
	for _, a := range s {
		if a.RoomName == e {
			return true
		}
	}
	return false
}
