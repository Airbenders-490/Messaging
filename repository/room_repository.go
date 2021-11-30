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
	GetRoom = `SELECT * FROM chat.room WHERE roomID=$1  ;`
	GetRoomFor = `SELECT * FROM chat.student_rooms WHERE student=$1  ;`
	editChatroomParticipant = `UPDATE chat.room SET Students=$2 WHERE RoomID=$1;`
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
	var room *domain.ChatRoom
	var studentText  []string
	var Allstudent []domain.Student

	err := r.dbSession.Query(GetRoom, roomID).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin, &room.Name, &studentText);
	if err!=nil {
		return nil,err ;
	}
	for _,ID := range studentText{
		var student *domain.Student
		student.ID=ID
		Allstudent = append(Allstudent,*student )

	}
	room.Students=Allstudent



	fmt.Println("Message entry: ", room.RoomID, room.Admin,room.Name)

	return room,err

}


func (r RoomRepository) GetRoomsFor(ctx context.Context, studentID string) (*domain.StudentChatRooms, error) {
	var stuentRoom *domain.StudentChatRooms
	var roomsID []string
	var rooms []domain.ChatRoom
	err := r.dbSession.Query(GetRoomFor, studentID).Consistency(gocql.One).Scan(&stuentRoom.Student, &roomsID);
	if err!=nil {
		return nil,err ;
	}
	for _,ID := range roomsID{
		var room *domain.ChatRoom
		room.RoomID=ID
		rooms = append(rooms,*room )

	}
	stuentRoom.Rooms=rooms
return stuentRoom , err
}

func (r RoomRepository) EditChatRoomParticipants(ctx context.Context, roomID string, student []domain.Student) error {
var studentID []string

	for _,s := range student{

		studentID = append(studentID,s.ID )

	}
	err := r.dbSession.Query(editChatroomParticipant, roomID, studentID).Consistency(gocql.One);
	if err!=nil {
		log.Println("unable to get the school name")
	}

	return  nil
}

func (r RoomRepository) AddRoomForParticipants(ctx context.Context, roomID string, student []domain.StudentChatRooms) error {

	panic("implement me")
}

func (r RoomRepository) RemoveRoomForParticipants(ctx context.Context, roomID string, student []domain.Student) error {
	panic("implement me")
}

func (r RoomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	panic("implement me")
}






