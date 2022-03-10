package repository

import (
	"chat/domain"
	"chat/messaging/repository/cassandra"
	"context"
	"github.com/gocql/gocql"
)

type RoomRepository struct {
	dbSession cassandra.SessionInterface
}

func NewRoomRepository(session cassandra.SessionInterface) *RoomRepository {
	return &RoomRepository{
		dbSession: session,
	}
}

const (
	// chat.room queries
	addParticipantToRoom      = `UPDATE chat.room SET students = students +? WHERE roomid=?;`
	deleteRoom                = `DELETE FROM chat.room WHERE roomid=?;`
	getRoom                   = `SELECT * FROM chat.room WHERE roomid=?;`
	getChatRoomsByClass       = `SELECT * FROM chat.room WHERE class=? ALLOW FILTERING;`
	removeParticipantFromRoom = `UPDATE chat.room SET students = students -? WHERE roomid=?;`
	saveRoom                  = `INSERT INTO chat.room (roomid, name, admin, students) VALUES (?,?,?,?);`

	// chat.student_rooms queries
	addRoomForParticipant    = `UPDATE chat.student_rooms SET rooms = rooms +? WHERE student=?;`
	getRoomsFor              = `SELECT * FROM chat.student_rooms WHERE student=?;`
	removeRoomForParticipant = `UPDATE chat.student_rooms SET rooms = rooms-? WHERE student=?;`
)

func (r RoomRepository) AddParticipantToRoom(ctx context.Context, userID string, roomID string) error {
	return r.dbSession.Query(addParticipantToRoom, [1]string{userID}, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	return r.dbSession.Query(deleteRoom, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) GetRoom(ctx context.Context, roomID string) (*domain.ChatRoom, error) {
	var room domain.ChatRoom
	var studentText []string
	var allStudents []domain.Student

	err := r.dbSession.Query(getRoom, roomID).WithContext(ctx).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin.ID, &room.Deleted, &room.Name, &studentText)
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

func (r RoomRepository) RemoveParticipantFromRoom(ctx context.Context, userID string, roomID string) error {
	return r.dbSession.Query(removeParticipantFromRoom, [1]string{userID}, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) SaveRoom(ctx context.Context, room *domain.ChatRoom) error {
	var userIDArr []string
	for _, student := range room.Students {
		userIDArr = append(userIDArr, student.ID)
	}
	return r.dbSession.Query(saveRoom, room.RoomID, room.Name, room.Admin.ID, userIDArr).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) AddRoomForParticipant(ctx context.Context, roomID string, userID string) error {
	return r.dbSession.Query(addRoomForParticipant, [1]string{roomID}, userID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) AddRoomForParticipants(ctx context.Context, roomID string, userIDs []string) error {
	for _, id := range userIDs {
		err := r.dbSession.Query(addRoomForParticipant, [1]string{roomID}, id).WithContext(ctx).Consistency(gocql.One).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r RoomRepository) GetRoomsFor(ctx context.Context, userID string) (*domain.StudentChatRooms, error) {
	var StudentRoom domain.StudentChatRooms
	var roomsID []string // used to unmarshall entire set to [0]th entry in string array
	var rooms []domain.ChatRoom
	err := r.dbSession.Query(getRoomsFor, userID).WithContext(ctx).Consistency(gocql.One).Scan(&StudentRoom.Student.ID, &roomsID)
	if err != nil {
		return nil, err
	}

	if len(roomsID) == 0 {
		return &StudentRoom, nil
	}
	for _, ID := range roomsID {
		var room domain.ChatRoom
		room.RoomID = ID
		rooms = append(rooms, room)
	}
	StudentRoom.Rooms = rooms

	return &StudentRoom, err
}

func (r RoomRepository) GetChatRoomsByClass(ctx context.Context, className string) ([]domain.ChatRoom, error) {
	retrievedChatRooms := make([]domain.ChatRoom,0)
	var scanner cassandra.ScannerInterface
	var studentIDs []string
	scanner = r.dbSession.Query(getChatRoomsByClass, className).WithContext(ctx).Consistency(gocql.One).Iter().Scanner()

	for scanner.Next() {
		var room domain.ChatRoom
		err := scanner.Scan(&room.RoomID, &room.Admin.ID, &room.Class, &room.Deleted, &room.Name, &studentIDs)

		if err != nil {
			return nil, err
		}

		for _, ID := range studentIDs {
			var student domain.Student
			student.ID = ID
			room.Students = append(room.Students, student)
		}
		retrievedChatRooms = append(retrievedChatRooms, room)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return retrievedChatRooms, nil
}

func (r RoomRepository) RemoveRoomForParticipant(ctx context.Context, roomID string, userID string) error {
	return r.dbSession.Query(removeRoomForParticipant, [1]string{roomID}, userID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) RemoveRoomForParticipants(ctx context.Context, roomID string, students []domain.Student) error {
	roomIDArr := [1]string{roomID}
	for _, student := range students {
		// Uses roomID Array bc cannot marshal string into set type
		err := r.dbSession.Query(removeRoomForParticipant, roomIDArr, student.ID).WithContext(ctx).Consistency(gocql.One).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r RoomRepository) SaveRoomAndAddRoomForAllParticipants(ctx context.Context, room *domain.ChatRoom) error {
	batch := r.dbSession.NewBatch(cassandra.BatchUnlogged).WithContext(ctx)

	var userIDArr []string
	for _, user := range room.Students {
		userIDArr = append(userIDArr, user.ID)
	}
	// SaveRoom for chat.room
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: saveRoom,
		Args: []interface{}{room.RoomID, room.Name, room.Admin.ID, userIDArr},
	})

	// AddRoomForAllParticipants for chat.student_rooms
	for _, user := range room.Students {
		batch.AddBatchEntry(&gocql.BatchEntry{
			Stmt: addRoomForParticipant,
			Args: []interface{}{[1]string{room.RoomID}, user.ID},
		})
	}
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) RemoveRoomForParticipantsAndDeleteRoom(ctx context.Context, room *domain.ChatRoom) error {
	batch := r.dbSession.NewBatch(cassandra.BatchUnlogged).WithContext(ctx)

	// RemoveRoomForParticipants for chat.student_rooms
	for _, user := range room.Students {
		batch.AddBatchEntry(&gocql.BatchEntry{
			Stmt: removeRoomForParticipant,
			Args: []interface{}{[1]string{room.RoomID}, user.ID},
		})
	}
	// DeleteRoom for chat.room
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: deleteRoom,
		Args: []interface{}{room.RoomID},
	})
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) AddParticipantToRoomAndAddRoomForParticipant(ctx context.Context, roomID string, userID string) error {
	batch := r.dbSession.NewBatch(cassandra.BatchUnlogged).WithContext(ctx)
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: addParticipantToRoom,
		Args: []interface{}{[1]string{userID}, roomID},
	})
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: addRoomForParticipant,
		Args: []interface{}{[1]string{roomID}, userID},
	})
	return r.dbSession.ExecuteBatch(batch)
}

func (r RoomRepository) RemoveParticipantFromRoomAndRemoveRoomForParticipant(ctx context.Context, roomID string, userID string) error {
	batch := r.dbSession.NewBatch(cassandra.BatchUnlogged).WithContext(ctx)
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: removeParticipantFromRoom,
		Args: []interface{}{[1]string{userID}, roomID},
	})
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: removeRoomForParticipant,
		Args: []interface{}{[1]string{roomID}, userID},
	})
	return r.dbSession.ExecuteBatch(batch)
}
