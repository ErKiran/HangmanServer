package events

type Room struct {
	RoomName   string `json:"roomName"`
	SecretWord string `json:"secretWord"`
}

type JoinRoom struct {
	NickName string `json:"nickName"`
	RoomName string `json:"roomName"`
}

type RoomsAndUser struct {
	RoomName string `json:"roomName"`
	NickName string `json:"nickName"`
}

var rooms []Room
var roomsAndUser []RoomsAndUser

func CreateRooms(msg Room) []Room {
	alreadyExists := contains(rooms, msg.RoomName)
	if !alreadyExists {
		rooms = append(rooms, Room{RoomName: msg.RoomName, SecretWord: msg.SecretWord})
	}
	return rooms
}

func GetRooms() []Room {
	return rooms
}

func DoRoomAlreadyExists(roomName string) bool {
	var isExists bool
	for _, name := range rooms {
		if name.RoomName == roomName {
			isExists = true
		}
	}
	return isExists
}

func JoinRooms(msg JoinRoom) []RoomsAndUser {
	roomsAndUser = append(roomsAndUser, RoomsAndUser{RoomName: msg.RoomName, NickName: msg.NickName})
	return roomsAndUser
}

func GetUserOfRooms(roomName string) []string {
	var users []string
	for _, room := range roomsAndUser {
		if room.RoomName == roomName {
			users = append(users, room.NickName)
		}
	}
	return users
}

func contains(s []Room, e string) bool {
	for _, a := range s {
		if a.RoomName == e {
			return true
		}
	}
	return false
}
