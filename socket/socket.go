package socket

import (
	"fmt"
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		next.ServeHTTP(w, r)
	})
}

func Socket() {

	var games []CreateGame
	var players []JoinGame

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
		fmt.Println("alreadyExists", alreadyExists)
		if !alreadyExists {
			games = append(games, CreateGame{RoomName: msg.RoomName, SecretWord: msg.SecretWord})
		}
	})

	server.OnEvent("/", "join-game", func(s socketio.Conn, msg JoinGame) {
		players = append(players, JoinGame{NickName: msg.NickName, RoomName: msg.RoomName})
		s.Emit("all-games", games)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	http.Handle("/socket.io/", corsMiddleware(server))

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
