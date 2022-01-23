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
	AddUserToRoom(roomID string, userID string) error
	RemoveUserFromRoom(roomID string, userID string) error
	GetChatRoomsFor( userID string) (*StudentChatRooms, error)
	DeleteRoom( userID string, roomID string) error
*/
type roomUseCase struct {
	rr      domain.RoomRepository
	sr      domain.StudentRepository
	timeout time.Duration
}

func NewRoomUseCase(rr domain.RoomRepository, sr domain.StudentRepository, t time.Duration) domain.RoomUseCase {
	return &roomUseCase{rr: rr, sr: sr, timeout: t}
}

// SaveRoom should add room to chat.room & chat.student_rooms for all participants
func (u *roomUseCase) SaveRoom(room *domain.ChatRoom) error {

	_, err := u.rr.GetRoom(room.RoomID)
	if err == nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", room.RoomID))
	}
	_, err = u.sr.GetStudent(room.Admin.ID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("The admin with ID  %s does not exist", room.Admin.ID))
	}

	err = u.rr.SaveRoomAndAddRoomForAllParticipants(room)
	if err != nil {
		return errors.NewInternalServerError("Unable to Save Room")
	}
	return err
}

// AddUserToRoom should add user to room in chat.room and add room to student in chat.student_rooms
func (u *roomUseCase) AddUserToRoom(roomID string, userID string) error {
	err := u.rr.AddParticipantToRoomAndAddRoomForParticipant(roomID, userID)
	if err != nil {
		return errors.NewInternalServerError("Unable to Add User to Room")
	}
	return err
}

// RemoveUserFromRoom should remove user from room in chat.room and remove room from user in chat.student_rooms
func (u *roomUseCase) RemoveUserFromRoom(roomID string, userID string) error {
	err := u.rr.RemoveParticipantFromRoomAndRemoveRoomForParticipant(roomID, userID)
	if err != nil {
		return errors.NewInternalServerError("Unable to Remove User from Room")
	}
	return err
}

// GetChatRoomsFor should get rooms for user in chat.student_rooms
func (u *roomUseCase) GetChatRoomsFor(userID string) (*domain.StudentChatRooms, error) {

	_, err := u.sr.GetStudent(userID)
	if err != nil {
		return nil, errors.NewConflictError(fmt.Sprintf("The Student with ID %s does not exist", userID))
	}
	var studentChatRooms *domain.StudentChatRooms

	studentChatRooms, err = u.rr.GetRoomsFor(userID)
	if err != nil {
		return nil, err
	}

	for _, Room := range studentChatRooms.Rooms {
		var room *domain.ChatRoom
		var student *domain.Student
		room, err = u.rr.GetRoom(Room.RoomID)

		student, err = u.sr.GetStudent(room.Admin.ID)
		room.Admin.LastName = student.LastName
		room.Admin.FirstName = student.FirstName

		Room.RoomID = room.RoomID
		Room.Name = room.Name
		// should we fill all students  ?? data
		Room.Students = room.Students
		Room.Deleted = room.Deleted
	}

	return studentChatRooms, nil
}

// DeleteRoom should delete room for all users in chat.student_rooms and delete room from chat.room
func (u *roomUseCase) DeleteRoom(userID string, roomID string) error {
	room, err := u.rr.GetRoom(roomID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does NOT exist", roomID))
	}

	if room.Admin.ID != userID {
		return errors.NewUnauthorizedError("Unauthorized to delete room, you are not Admin")
	}

	err = u.rr.RemoveRoomForParticipantsAndDeleteRoom(room)
	if err != nil {
		return errors.NewInternalServerError("Unable to Remove User from Room")
	}
	return err
}
