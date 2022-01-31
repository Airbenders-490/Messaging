package usecase

import (
	"chat/domain"
	"chat/utils/errors"
	"context"
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
func (u *roomUseCase) SaveRoom(ctx context.Context, room *domain.ChatRoom) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	_, err := u.rr.GetRoom(ctx, room.RoomID)
	if err == nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s already exists", room.RoomID))
	}

	for _, participant := range room.Students {
		_, err = u.sr.GetStudent(ctx, participant.ID)
		if err != nil {
			return errors.NewConflictError(fmt.Sprintf("The participant with ID %s does not exist", participant.ID))
		}
	}

	return u.rr.SaveRoomAndAddRoomForAllParticipants(ctx, room)
}

// AddUserToRoom should add user to room in chat.room and add room to student in chat.student_rooms
func (u *roomUseCase) AddUserToRoom(ctx context.Context, roomID string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.rr.AddParticipantToRoomAndAddRoomForParticipant(ctx, roomID, userID)
}

// RemoveUserFromRoom should remove user from room in chat.room and remove room from user in chat.student_rooms
func (u *roomUseCase) RemoveUserFromRoom(ctx context.Context, roomID string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.rr.RemoveParticipantFromRoomAndRemoveRoomForParticipant(ctx, roomID, userID)
}

// GetChatRoomsFor should get rooms for user in chat.student_rooms
func (u *roomUseCase) GetChatRoomsFor(ctx context.Context, userID string) (*domain.StudentChatRooms, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	_, err := u.sr.GetStudent(ctx, userID)
	if err != nil {
		return nil, errors.NewConflictError(fmt.Sprintf("The Student with ID %s does not exist", userID))
	}
	var studentChatRooms *domain.StudentChatRooms

	studentChatRooms, err = u.rr.GetRoomsFor(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, Room := range studentChatRooms.Rooms {
		var room *domain.ChatRoom
		var student *domain.Student
		room, err = u.rr.GetRoom(ctx, Room.RoomID)
		if err != nil {
			return nil, err
		}
		student, err = u.sr.GetStudent(ctx, room.Admin.ID)
		if err != nil {
			return nil, err
		}

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
func (u *roomUseCase) DeleteRoom(ctx context.Context, userID string, roomID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	room, err := u.rr.GetRoom(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Admin.ID != userID {
		return errors.NewUnauthorizedError("Unauthorized to delete room, you are not Admin")
	}

	return u.rr.RemoveRoomForParticipantsAndDeleteRoom(ctx, room)
}
