package domain

import (
	"context"
	"time"
)

// ChatRoom struct
type ChatRoom struct {
	RoomID string
	Name string
	Admin Student
	Deleted   time.Time
	Students []Student
}

// StudentChatRooms struct
type StudentChatRooms struct {
	Student Student
	Rooms []ChatRoom
}

// RoomRepository interface implements the contract as descirbed aboved each method
type RoomRepository interface {
	SaveRoom( room *ChatRoom) error
	GetRoom(roomID string) (*ChatRoom, error)
	GetRoomsFor( studentID string) (*StudentChatRooms, error)
	// EditChatRoomParticipants deals with chat.room and changes students there
	EditChatRoomParticipants( roomID string, student []Student) error
	// AddRoomForParticipants deals with chat.student_rooms and adds the chatroom to each student's list
	AddRoomForParticipants( roomID string, student []Student) error
	// RemoveRoomForParticipants deals with chat.student_rooms and removes the chatroom from each student's list
	RemoveRoomForParticipants( roomID string, student []Student) error
	DeleteRoom(roomID string) error
}

// RoomUseCase interface implements the contract as described above each method
type RoomUseCase interface {
	// SaveRoom needs to save not just the room, but also add the chatroom for all the students
	SaveRoom(ctx context.Context, room *ChatRoom) error
	GetChatRoomsFor(ctx context.Context, studentID string) (*StudentChatRooms, error)
	// EditChatRoomParticipants this deals with not only changing the participants in the chat.rooms table, but also
	// manipulating the chat.student_rooms table to ensure updated room list
	EditChatRoomParticipants(ctx context.Context, room *ChatRoom) error
	// DeleteRoom Ensure the user deleting is the admin
	DeleteRoom(ctx context.Context, userID string, roomID string) error
}


