package repository

import (
	"chat/domain"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

/*
// ChatRoom struct
type ChatRoom struct {
	RoomID string
	Name string
	Admin Student
	Students []Student
}

// StudentChatRooms struct
type StudentChatRooms struct {
	Student Student
	Rooms []ChatRoom
}

type RoomRepository interface {
	SaveRoom(ctx context.Context, room *ChatRoom) error
	GetRoom(ctx context.Context, roomID string) (*ChatRoom, error)
	GetRoomsFor(ctx context.Context, studentID string) (*StudentChatRooms, error)
	// EditChatRoomParticipants deals with chat.room and changes students there
	EditChatRoomParticipants(ctx context.Context, roomID string, student []Student) error
	// AddRoomForParticipants deals with chat.student_rooms and adds the chatroom to each student's list
	AddRoomForParticipants(ctx context.Context, roomID string, student []Student) error
	// RemoveRoomForParticipants deals with chat.student_rooms and removes the chatroom from each student's list
	RemoveRoomForParticipants(ctx context.Context, roomID string, student []Student) error
	DeleteRoom(ctx context.Context, roomID string) error
}*/

type RoomRepository struct {
	dbSession *gocql.Session
}
const (
	SaveRoom = `INSERT INTO chat.room (roomID, NAME, Admin, students) VALUES ($1, $2, $3, $4);`
)


func (r RoomRepository) SaveRoom(ctx context.Context, room *domain.ChatRoom) error {
	var students []string
	for _, element := range room.Students{
		students = append(students, element.ID)
	}
	 err := r.dbSession.Query(SaveRoom, room.RoomID, room.Admin,room.Name,students).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin, &room.Name, &students);

	if err!=nil {
		log.Fatal(err)
	}


	fmt.Println("Message entry: ", room.RoomID, room.Admin,room.Name,room.Students)

	return err

}

func (r RoomRepository) GetRoom(ctx context.Context, roomID string) (*domain.ChatRoom, error) {
	panic("implement me")
}

func (r RoomRepository) GetRoomsFor(ctx context.Context, studentID string) (*domain.StudentChatRooms, error) {
	panic("implement me")
}

func (r RoomRepository) EditChatRoomParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r RoomRepository) AddRoomForParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r RoomRepository) RemoveRoomForParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r RoomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	panic("implement me")
}

func NewRoomRepository(session *gocql.Session) domain.RoomRepository {
	return &RoomRepository{
		dbSession: session,
	}
}



