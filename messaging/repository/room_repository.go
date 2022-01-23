package repository

import (
	"chat/domain"
	"context"
	"github.com/gocql/gocql"
	"strings"
)

type RoomRepository struct {
	dbSession *gocql.Session
}

func NewRoomRepository(session *gocql.Session) *RoomRepository {
	return &RoomRepository{
		dbSession: session,
	}
}

const (
	// chat.room queries
	addParticipantToRoom      = `UPDATE chat.room SET students = students +? WHERE roomid=?;`
	deleteRoom                = `DELETE FROM chat.room WHERE roomid=?;`
	getRoom                   = `SELECT * FROM chat.room WHERE roomid=? LIMIT 1;`
	removeParticipantFromRoom = `UPDATE chat.room SET students = students -? WHERE roomid=?;`
	saveRoom                  = `INSERT INTO chat.room (roomid, name, admin, students) VALUES (?,?,?,?);`

	// chat.student_rooms queries
	addRoomForParticipant    = `UPDATE chat.student_rooms SET rooms = rooms +? WHERE student=?;`
	getRoomsFor              = `SELECT * FROM chat.student_rooms WHERE student=?;`
	removeRoomForParticipant = `UPDATE chat.student_rooms SET rooms = rooms-? WHERE student=?;`
)

func (r RoomRepository) AddParticipantToRoom(studentID string, roomID string) error {
	return r.dbSession.Query(addParticipantToRoom, [1]string{studentID}, roomID).Consistency(gocql.One).Exec()
}

func (r RoomRepository) DeleteRoom(roomID string) error {
	return r.dbSession.Query(deleteRoom, roomID).Consistency(gocql.One).Exec()
}

func (r RoomRepository) GetRoom(roomID string) (*domain.ChatRoom, error) {
	var room domain.ChatRoom
	var studentText []string
	var allStudents []domain.Student

	err := r.dbSession.Query(getRoom, roomID).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin.ID, &room.Deleted, &room.Name, &studentText)
	if err != nil {
		return nil, err
	}

	for _, ID := range studentText {
		var student domain.Student
		student.ID = ID
		allStudents = append(allStudents, student)
	}
	room.Students = allStudents

	return &room, err
}

func (r RoomRepository) RemoveParticipantFromRoom(studentID string, roomID string) error {
	return r.dbSession.Query(removeParticipantFromRoom, [1]string{studentID}, roomID).Consistency(gocql.One).Exec()
}

func (r RoomRepository) SaveRoom(room *domain.ChatRoom) error {
	var studentIDArr []string
	for _, student := range room.Students {
		studentIDArr = append(studentIDArr, student.ID)
	}
	return r.dbSession.Query(saveRoom, room.RoomID, room.Name, room.Admin.ID, studentIDArr).Consistency(gocql.One).Exec()
}

func (r RoomRepository) AddRoomForParticipant(roomID string, studentID string) error {
	return r.dbSession.Query(addRoomForParticipant, [1]string{roomID}, studentID).Consistency(gocql.One).Exec()
}

func (r RoomRepository) AddRoomForParticipants(roomID string, studentIDs []string) error {
	for _, id := range studentIDs {
		err := r.dbSession.Query(addRoomForParticipant, [1]string{roomID}, id).Consistency(gocql.One).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r RoomRepository) GetRoomsFor(studentID string) (*domain.StudentChatRooms, error) {
	var StudentRoom domain.StudentChatRooms
	var roomsID []string // used to unmarshall entire set to [0]th entry in string array
	var rooms []domain.ChatRoom
	err := r.dbSession.Query(getRoomsFor, studentID).Consistency(gocql.One).Scan(&StudentRoom.Student.ID, &roomsID)
	if err != nil {
		return nil, err
	}

	roomsIDSplit := strings.Split(roomsID[0], ",")
	for _, ID := range roomsIDSplit {
		ID = strings.ReplaceAll(ID, " ", "")
		var room domain.ChatRoom
		room.RoomID = ID
		rooms = append(rooms, room)
	}
	StudentRoom.Rooms = rooms

	return &StudentRoom, err
}

func (r RoomRepository) RemoveRoomForParticipant(roomID string, studentID string) error {
	return r.dbSession.Query(removeRoomForParticipant, [1]string{roomID}, studentID).Consistency(gocql.One).Exec()
}

func (r RoomRepository) RemoveRoomForParticipants(roomID string, students []domain.Student) error {
	roomIDArr := [1]string{roomID}
	for _, student := range students {
		// Uses roomID Array bc cannot marshal string into set type
		err := r.dbSession.Query(removeRoomForParticipant, roomIDArr, student.ID).Consistency(gocql.One).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r RoomRepository) SaveRoomAndAddRoomForAllParticipants(room *domain.ChatRoom) error {
	batch := r.dbSession.NewBatch(gocql.UnloggedBatch).WithContext(context.Background())

	var userIDArr []string
	for _, user := range room.Students {
		userIDArr = append(userIDArr, user.ID)
	}
	// SaveRoom for chat.room
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: saveRoom,
		Args: []interface{}{room.RoomID, room.Name, room.Admin.ID, userIDArr},
	})

	// AddRoomForAllParticipants for chat.student_rooms
	for _, user := range room.Students {
		batch.Entries = append(batch.Entries, gocql.BatchEntry{
			Stmt: addRoomForParticipant,
			Args: []interface{}{[1]string{room.RoomID}, user.ID},
		})
	}
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) RemoveRoomForParticipantsAndDeleteRoom(room *domain.ChatRoom) error {
	batch := r.dbSession.NewBatch(gocql.UnloggedBatch).WithContext(context.Background())

	// RemoveRoomForParticipants for chat.student_rooms
	for _, user := range room.Students {
		batch.Entries = append(batch.Entries, gocql.BatchEntry{
			Stmt: removeRoomForParticipant,
			Args: []interface{}{[1]string{room.RoomID}, user.ID},
		})
	}
	// DeleteRoom for chat.room
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: deleteRoom,
		Args: []interface{}{room.RoomID},
	})
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) AddParticipantToRoomAndAddRoomForParticipant(roomID string, studentID string) error {
	batch := r.dbSession.NewBatch(gocql.UnloggedBatch).WithContext(context.Background())
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: addParticipantToRoom,
		Args: []interface{}{[1]string{studentID}, roomID},
	})
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: addRoomForParticipant,
		Args: []interface{}{[1]string{roomID}, studentID},
	})
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) RemoveParticipantFromRoomAndRemoveRoomForParticipant(roomID string, studentID string) error {
	batch := r.dbSession.NewBatch(gocql.UnloggedBatch).WithContext(context.Background())
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: removeParticipantFromRoom,
		Args: []interface{}{[1]string{studentID}, roomID},
	})
	batch.Entries = append(batch.Entries, gocql.BatchEntry{
		Stmt: removeRoomForParticipant,
		Args: []interface{}{[1]string{roomID}, studentID},
	})
	return r.dbSession.ExecuteBatch(batch)
}
