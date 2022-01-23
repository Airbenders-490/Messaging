package domain

import (
	"time"
)

// ChatRoom struct
type ChatRoom struct {
	RoomID   string    `json:"room_id"`
	Name     string    `json:"name"`
	Admin    Student   `json:"admin"`
	Deleted  time.Time `json:"deleted"`
	Students []Student `json:"students"`
}

// StudentChatRooms struct
type StudentChatRooms struct {
	Student Student
	Rooms   []ChatRoom
}

// RoomRepository interface implements the contract as descirbed aboved each method
type RoomRepository interface {

	// chat.room methods
	AddParticipantToRoom(userID string, roomID string) error
	DeleteRoom(roomID string) error
	GetRoom(roomID string) (*ChatRoom, error)
	RemoveParticipantFromRoom(userID string, roomID string) error
	SaveRoom(room *ChatRoom) error

	// chat.student_rooms methods
	AddRoomForParticipant(roomID string, userID string) error
	// AddRoomForParticipants deals with chat.student_rooms and adds the chatroom to each student's list
	AddRoomForParticipants(roomID string, userIDs []string) error
	GetRoomsFor(userID string) (*StudentChatRooms, error)
	RemoveRoomForParticipant(roomID string, userID string) error
	// RemoveRoomForParticipants deals with chat.student_rooms and removes the chatroom from each student's list
	RemoveRoomForParticipants(roomID string, users []Student) error

	// batch
	SaveRoomAndAddRoomForAllParticipants(room *ChatRoom) error
	RemoveRoomForParticipantsAndDeleteRoom(room *ChatRoom) error
	AddParticipantToRoomAndAddRoomForParticipant(roomID string, userID string) error
	RemoveParticipantFromRoomAndRemoveRoomForParticipant(roomID string, userID string) error
}

// RoomUseCase interface implements the contract as described above each method
type RoomUseCase interface {
	// SaveRoom needs to save not just the room, but also add the chatroom for all the students
	SaveRoom(room *ChatRoom) error
	AddUserToRoom(roomID string, userID string) error
	RemoveUserFromRoom(roomID string, userID string) error
	GetChatRoomsFor(userID string) (*StudentChatRooms, error)
	// DeleteRoom Ensure the user deleting is the admin
	DeleteRoom(userID string, roomID string) error
}
