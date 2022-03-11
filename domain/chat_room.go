package domain

import (
	"context"
	"time"
)

// ChatRoom struct
type ChatRoom struct {
	RoomID   string    `json:"room_id"`
	Name     string    `json:"name"`
	Admin    Student   `json:"admin"`
	Deleted  time.Time `json:"deleted"`
	Students []Student `json:"students"`
	Class    string    `json:"class"`
}

// StudentChatRooms struct
type StudentChatRooms struct {
	Student Student
	Rooms   []ChatRoom
}

// RoomRepository interface implements the contract as descirbed aboved each method
type RoomRepository interface {

	// chat.room methods
	AddParticipantToRoom(ctx context.Context, userID string, roomID string) error
	DeleteRoom(ctx context.Context, roomID string) error
	GetRoom(ctx context.Context, roomID string) (*ChatRoom, error)
	GetChatRoomsByClass(ctx context.Context, className string) ([]ChatRoom, error)
	RemoveParticipantFromRoom(ctx context.Context, userID string, roomID string) error
	SaveRoom(ctx context.Context, room *ChatRoom) error

	// chat.student_rooms methods
	AddRoomForParticipant(ctx context.Context, roomID string, userID string) error
	// AddRoomForParticipants deals with chat.student_rooms and adds the chatroom to each student's list
	AddRoomForParticipants(ctx context.Context, roomID string, userIDs []string) error
	GetRoomsFor(ctx context.Context, userID string) (*StudentChatRooms, error)
	RemoveRoomForParticipant(ctx context.Context, roomID string, userID string) error
	// RemoveRoomForParticipants deals with chat.student_rooms and removes the chatroom from each student's list
	RemoveRoomForParticipants(ctx context.Context, roomID string, users []Student) error

	// batch
	SaveRoomAndAddRoomForAllParticipants(ctx context.Context, room *ChatRoom) error
	RemoveRoomForParticipantsAndDeleteRoom(ctx context.Context, room *ChatRoom) error
	AddParticipantToRoomAndAddRoomForParticipant(ctx context.Context, roomID string, userID string) error
	RemoveParticipantFromRoomAndRemoveRoomForParticipant(ctx context.Context, roomID string, userID string) error
}

// RoomUseCase interface implements the contract as described above each method
type RoomUseCase interface {
	// SaveRoom needs to save not just the room, but also add the chatroom for all the students
	SaveRoom(ctx context.Context, room *ChatRoom) error
	AddUserToRoom(ctx context.Context, roomID string, userID string) error
	RemoveUserFromRoom(ctx context.Context, roomID string, userID string, loggedID string) error
	GetChatRoomsByClass(ctx context.Context, className string) ([]ChatRoom, error)
	GetChatRoomsFor(ctx context.Context, userID string) (*StudentChatRooms, error)
	// DeleteRoom Ensure the user deleting is the admin
	DeleteRoom(ctx context.Context, userID string, roomID string) error
}
