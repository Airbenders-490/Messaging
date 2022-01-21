package usecase

import (
	"chat/domain"
	"chat/utils/errors"
	"fmt"
	_ "reflect"
	"time"
)
/*

	SaveRoom( room *ChatRoom) error
	AddUserToRoom(roomID string, studentID string) error
	RemoveUserFromRoom(roomID string, studentID string) error
	GetChatRoomsFor( studentID string) (*StudentChatRooms, error)
	DeleteRoom( userID string, roomID string) error
*/
type roomUseCase struct {
	rr      domain.RoomRepository
	sr 		domain.StudentRepository
	timeout time.Duration
}

func NewRoomUseCase(rr domain.RoomRepository, sr domain.StudentRepository, t time.Duration) domain.RoomUseCase {
	return &roomUseCase{rr: rr, sr: sr, timeout: t}
}

// SaveRoom should add room to chat.room & chat.student_rooms for all participants
func (u *roomUseCase) SaveRoom(room *domain.ChatRoom) error {

	_,err := u.rr.GetRoom(room.RoomID)
	if err==nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", room.RoomID))
	}
	_,err=u.sr.GetStudent(room.Admin.ID)
	if err!=nil {
		return errors.NewConflictError(fmt.Sprintf("The admin with ID  %s does not exist", room.Admin.ID))
	}
	err = u.rr.SaveRoom( room )
	if err!=nil {
		return err
	}

	var studentIDs []string
	for _,student := range room.Students {
		studentIDs = append(studentIDs, student.ID)
	}
	err = u.rr.AddRoomForParticipants(room.RoomID, studentIDs)
	if err!=nil {
		return err
	}

   	return err
}

// AddUserToRoom should add user to room in chat.room and add room to student in chat.student_rooms
func (u *roomUseCase) AddUserToRoom(roomID string, studentID string) error {

	err := u.rr.AddParticipantToRoom(studentID, roomID)
	if err!=nil {
		return errors.NewConflictError("Unable to Add Participant to Room")
	}

	err = u.rr.AddRoomForParticipant(roomID, studentID)
	if err!=nil {
		return errors.NewConflictError("Unable to Add Room for Participant")
	}
	return nil
}

// RemoveUserFromRoom should remove user from room in chat.room and remove room from user in chat.student_rooms
func (u *roomUseCase) RemoveUserFromRoom(roomID string, studentID string) error {

	err := u.rr.RemoveParticipantFromRoom(studentID, roomID)
	if err!=nil {
		return errors.NewConflictError("Unable to Remove Participant from Room")
	}

	err = u.rr.RemoveRoomForParticipant(roomID, studentID)
	if err!=nil {
		return errors.NewConflictError("Unable to Remove Room for Participant")
	}
	return nil
}

// GetChatRoomsFor should get rooms for user in chat.student_rooms
func (u *roomUseCase) GetChatRoomsFor(studentID string) (*domain.StudentChatRooms, error) {

	_,err:=u.sr.GetStudent(studentID)
	if err!=nil {
		return nil,errors.NewConflictError(fmt.Sprintf("The Student with ID %s does not exist", studentID))
	}
	var studentChatRooms *domain.StudentChatRooms

	studentChatRooms,err=u.rr.GetRoomsFor(studentID)
	if err!=nil {
		return nil, err
	}

	for _,Room := range studentChatRooms.Rooms{
		var room *domain.ChatRoom
		var student *domain.Student
		room,err=u.rr.GetRoom(Room.RoomID)

		student,err=u.sr.GetStudent(room.Admin.ID)
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

// DeleteRoom should delete room for all users in chat.student_rooms and delete room from chat.room
func (u *roomUseCase) DeleteRoom( userID string, roomID string) error {
	room,err := u.rr.GetRoom(roomID)
	if err!=nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does NOT exist", roomID))
	}

	if room.Admin.ID != userID {
		return errors.NewConflictError("Unauthorized to delete room, you are not Admin")
	}

 	err = u.rr.RemoveRoomForParticipants(roomID, room.Students)
	if err!=nil {
		return err
	}

	err = u.rr.DeleteRoom(roomID)
	if err!=nil {
		return err
	}

	return nil
}





