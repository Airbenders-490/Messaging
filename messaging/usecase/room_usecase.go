package usecase

import (
	"chat/domain"
	"chat/utils/errors"
	"fmt"
	"log"
	_ "reflect"
	"time"
)
/*

	SaveRoom( room *ChatRoom) error
	GetChatRoomsFor( studentID string) (*StudentChatRooms, error)
	EditChatRoomParticipants( room *ChatRoom) error
	DeleteRoom( userID string, roomID string) error
*/
type roomUseCase struct {
	rr      domain.RoomRepository
	st 		domain.StudentRepository
	timeout time.Duration
}

func (u *roomUseCase) SaveRoom(room *domain.ChatRoom) error {

	_,err := u.rr.GetRoom(room.RoomID)
	if err==nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", room.RoomID))
	}
	_,err=u.st.GetStudent(room.Admin.ID)
	if err==nil {
		return errors.NewConflictError(fmt.Sprintf("The admin with ID  %s does not exist", room.Admin.ID))
	}
	err = u.rr.SaveRoom( room )
	if err!=nil {
		return err
	}

   return err
}

func (u *roomUseCase)  GetStudentChatRoomsFor( studentID string) (*domain.StudentChatRooms, error){


	_,err:=u.st.GetStudent(studentID)
	if err==nil {
		return nil,errors.NewConflictError(fmt.Sprintf("The Student with ID  %s does not exist", studentID))
	}
	var studentChatRooms *domain.StudentChatRooms

	studentChatRooms,err=u.rr.GetRoomsFor(studentID)
if err!=nil {
	log.Fatal(err)
	return nil , err
}
	for _,Room := range studentChatRooms.Rooms{
		var room *domain.ChatRoom
		var student *domain.Student
		room,err=u.rr.GetRoom(Room.RoomID)

		student,err=u.st.GetStudent(room.Admin.ID)
		room.Admin.LastName=student.LastName
		room.Admin.FirstName=student.FirstName

		Room.RoomID=room.RoomID ;
		Room.Name=room.Name ;
		// should we fill all students  ?? data
		Room.Students=room.Students ;
		Room.Deleted=room.Deleted ;

	}

	return studentChatRooms,nil
}

func (u *roomUseCase) EditChatRoomParticipants( room *domain.ChatRoom) error{

	_,err:=u.rr.GetRoom(room.RoomID)
	if err==nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", room.RoomID))
	}

		err =u.rr.EditChatRoomParticipants( room.RoomID , room.Students)
if err!=nil {

	log.Fatal(err)
	return  err
}


	return nil
}

func (u *roomUseCase) DeleteRoom( userID string, roomID string) error {
	room,err := u.rr.GetRoom(roomID)
	if err!=nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", roomID))
	}

 err = u.rr.RemoveRoomForParticipants(roomID, room.Students)
if err!=nil {
	log.Fatal(err)
	return err
}
	return nil
}





