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

	student, err := u.sr.GetStudent(ctx, room.Admin.ID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("User %s does not exist", room.Admin.ID))
	}
	room.Admin.LastName = student.LastName
	room.Admin.FirstName = student.FirstName

	room.Students = append(room.Students, room.Admin)
	for _, participant := range room.Students {
		_, err = u.sr.GetStudent(ctx, participant.ID)
		if err != nil {
			return errors.NewConflictError(fmt.Sprintf("The participant with ID %s does not exist", participant.ID))
		}
	}

	return u.rr.SaveRoomAndAddRoomForAllParticipants(ctx, room)
}

// AddUserToRoom should add user to room in chat.room and add room to student in chat.student_rooms
func (u *roomUseCase) AddUserToRoom(ctx context.Context, roomID string, userID string, loggedID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	room, err := u.rr.GetRoom(ctx, roomID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does not exist", roomID))
	}

	if room.Admin.ID != loggedID {
		return errors.NewUnauthorizedError("Unauthorized, you cannot add a user unless you are the admin")
	}

	isPendingFalseCount := 0
	for _, student := range room.Students {
		if (!student.IsPending) {
			isPendingFalseCount++
		}
	}

	if isPendingFalseCount < room.MaxParticipants {
		return u.rr.AddParticipantToRoomAndAddRoomForParticipant(ctx, roomID, userID)
	}
	return errors.NewConflictError("Room is full")
}

// RemoveUserFromRoom should remove user from room in chat.room and remove room from user in chat.student_rooms
func (u *roomUseCase) RemoveUserFromRoom(ctx context.Context, roomID string, userID string, loggedID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	room, err := u.rr.GetRoom(ctx, roomID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does not exist", roomID))
	}

	if room.Admin.ID != loggedID && userID != loggedID  {
		return errors.NewUnauthorizedError("Unauthorized, you cannot remove someone else unless you are the admin")
	}


	return u.rr.RemoveParticipantFromRoomAndRemoveRoomForParticipant(ctx, roomID, userID)
}

// GetChatRoomsFor should get rooms for user in chat.student_rooms
func (u *roomUseCase) GetChatRoomsFor(ctx context.Context, userID string) (*domain.StudentChatRooms, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	student , err := u.sr.GetStudent(ctx, userID)
	if err != nil {
		return nil, errors.NewConflictError(fmt.Sprintf("The Student with ID %s does not exist", userID))
	}

	studentChatRooms, err := u.rr.GetRoomsFor(ctx, userID)
	if err != nil {
		return nil, err
	}
	studentChatRooms.Student.FirstName = student.FirstName
	studentChatRooms.Student.LastName = student.LastName

	for i := range studentChatRooms.Rooms {
		var room *domain.ChatRoom
		var student *domain.Student
		room, err = u.rr.GetRoom(ctx, studentChatRooms.Rooms[i].RoomID)
		if err != nil {
			return nil, err
		}
		student, err = u.sr.GetStudent(ctx, room.Admin.ID)
		if err != nil {
			return nil, err
		}
		room.Admin.LastName = student.LastName
		room.Admin.FirstName = student.FirstName

		studentChatRooms.Rooms[i].RoomID = room.RoomID
		studentChatRooms.Rooms[i].Admin = room.Admin
		studentChatRooms.Rooms[i].Class = room.Class
		studentChatRooms.Rooms[i].Name = room.Name
		studentChatRooms.Rooms[i].Students = room.Students
		studentChatRooms.Rooms[i].Deleted = room.Deleted
		studentChatRooms.Rooms[i].MaxParticipants = room.MaxParticipants

		for j := range studentChatRooms.Rooms[i].Students {
			student, err = u.sr.GetStudent(ctx, studentChatRooms.Rooms[i].Students[j].ID)
			if err != nil {
				return nil, err
			}
			studentChatRooms.Rooms[i].Students[j].FirstName = student.FirstName
			studentChatRooms.Rooms[i].Students[j].LastName = student.LastName
		}
	}
	return studentChatRooms, nil
}

// GetChatRoomsByClass should get rooms in chat.room by className, returns empty list if no rooms found
func (u *roomUseCase) GetChatRoomsByClass(ctx context.Context, className string) ([]domain.ChatRoom, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	rooms, err := u.rr.GetChatRoomsByClass(ctx, className)
	if err != nil {
		return nil, err
	}

	for i := range rooms {
		var r *domain.ChatRoom
		var student *domain.Student
		r, err = u.rr.GetRoom(ctx, rooms[i].RoomID)
		if err != nil {
			return nil, err
		}
		student, err = u.sr.GetStudent(ctx, r.Admin.ID)
		if err != nil {
			return nil, err
		}
		r.Admin.LastName = student.LastName
		r.Admin.FirstName = student.FirstName

		rooms[i].RoomID = r.RoomID
		rooms[i].Name = r.Name
		rooms[i].Admin = r.Admin
		rooms[i].Class = r.Class
		rooms[i].Students = r.Students
		rooms[i].Deleted = r.Deleted

		for j := range rooms[i].Students {
			student, err = u.sr.GetStudent(ctx, rooms[i].Students[j].ID)
			if err != nil {
				return nil, err
			}
			rooms[i].Students[j].FirstName = student.FirstName
			rooms[i].Students[j].LastName = student.LastName
		}
	}
	return rooms, nil
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
