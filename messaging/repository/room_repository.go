package repository

import (
	"chat/domain"
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

func NewRoomRepository(session *gocql.Session) *RoomRepository {
	return &RoomRepository{
		dbSession: session,
	}
}
const (
	SaveRoom = `INSERT INTO chat.room (roomID, NAME, Admin, students) VALUES (?, ?, ?,?);`
	GetRoom = `SELECT * FROM chat.room WHERE roomID=? LIMIT 1;`
	GetRoomFor = `SELECT * FROM chat.student_rooms WHERE student=? LIMIT 1 ;`
	editChatroomParticipant = `UPDATE chat.room  SET Students=? WHERE RoomID=?;`

	addRoomForParticipant = `UPDATE chat.student_rooms  SET rooms = rooms +? WHERE student=? ;`
	CreateStudentRooms=`INSERT INTO chat.student_rooms (student) VALUES (?);`
	RemoveRoomForParticipant = `UPDATE chat.student_rooms  SET rooms = rooms-? WHERE student=? ;`
	DeleteRoom = `DELETE FROM chat.room   WHERE roomID=? ;`

)


func (r RoomRepository) SaveRoom( room *domain.ChatRoom) error {
	var students []string
	for _, element := range room.Students{
		students = append(students, element.ID)
	}

	 err := r.dbSession.Query(SaveRoom, room.RoomID, room.Admin.ID,room.Name,students).Consistency(gocql.One).Exec();

	if err!=nil {
	return err ;
	}


	fmt.Println("Message entry: ", room.RoomID, room.Admin,room.Name,room.Students)

	return err

}


func (r RoomRepository) GetRoom( roomID string) (*domain.ChatRoom, error) {
	var room domain.ChatRoom

	var studentText  []string
	var AllStudent []domain.Student

	err := r.dbSession.Query(GetRoom, roomID).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin.ID,&room.Deleted, &room.Name, &studentText);


	if err!=nil {
		return nil,err ;
	}

	for _,ID := range studentText{
		var student domain.Student
		student.ID=ID
		AllStudent = append(AllStudent,student )

	}
	room.Students= AllStudent





	return &room,err

}


func (r RoomRepository) EditChatRoomParticipants( roomID string, student []domain.Student) error {
var studentID []string

	for _,s := range student{

		studentID = append(studentID,s.ID )

	}
	err := r.dbSession.Query(editChatroomParticipant, studentID, roomID).Consistency(gocql.One).Exec();
	if err!=nil {
		return err ;
	}

	return  nil
}

func (r RoomRepository) GetRoomsFor( studentID string) (*domain.StudentChatRooms, error) {
	var StudentRoom domain.StudentChatRooms
	var roomsID []string
	var rooms []domain.ChatRoom
	err := r.dbSession.Query(GetRoomFor, studentID).Consistency(gocql.One).Scan(&StudentRoom.Student.ID, &roomsID);
	if err!=nil {
		return nil,err ;
	}
	for _,ID := range roomsID{
		var room domain.ChatRoom
		room.RoomID=ID
		rooms = append(rooms,room )

	}
	StudentRoom.Rooms=rooms

	return &StudentRoom , err
}



func (r RoomRepository) AddRoomForParticipants( roomID string, students []domain.Student) error {
	var rooms []string
	rooms=append(rooms,roomID)

	for _,student:= range students {

		err := r.dbSession.Query(addRoomForParticipant, rooms,student.ID).Consistency(gocql.One).Exec();
		if err!=nil {
			return err ;
		}
	}

	return  nil

}

func (r RoomRepository) CreateStudentRooms( studnet domain.Student) error {



		err := r.dbSession.Query(CreateStudentRooms, studnet.ID).Consistency(gocql.One).Exec();
		if err!=nil {
			return err ;
		}


	return  nil

}



func (r RoomRepository) RemoveRoomForParticipants( roomID string, students []domain.Student) error {
	var rooms []string
	rooms=append(rooms, roomID)
	for _,student:= range students {

		err := r.dbSession.Query(RemoveRoomForParticipant, rooms,student.ID).Exec();
		if err!=nil {
			return err ;
		}

	}

	return  nil

}

func (r RoomRepository) DeleteRoom( roomID string) error {
	err := r.dbSession.Query(DeleteRoom, roomID).Exec();
	if err!=nil {
		return err ;
	}

	return nil
}







