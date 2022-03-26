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
	deleteRoom                    = `DELETE FROM chat.room WHERE roomid=?;`
	getRoom                       = `SELECT * FROM chat.room WHERE roomid=?;`
	getChatRoomsByClass           = `SELECT * FROM chat.room WHERE class=? ALLOW FILTERING;`
	removeParticipantFromRoom     = `DELETE students[?] FROM chat.room WHERE roomid = ?;`
	saveRoom                      = `INSERT INTO chat.room (roomid, name, admin, students, class, maxParticipants) VALUES (?,?,?,?,?,?);`
	updateParticipantPendingState = `UPDATE chat.room SET students[?] = ?  WHERE roomid = ?;`

	// chat.student_rooms queries
	addRoomForParticipant    = `UPDATE chat.student_rooms SET rooms = rooms +? WHERE student=?;`
	getRoomsFor              = `SELECT * FROM chat.student_rooms WHERE student=?;`
	removeRoomForParticipant = `UPDATE chat.student_rooms SET rooms = rooms-? WHERE student=?;`
)

func (r RoomRepository) UpdateParticipantPendingState(ctx context.Context, roomID string, userID string, isPending bool) error {
	return r.dbSession.Query(updateParticipantPendingState, userID, isPending, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) DeleteRoom(ctx context.Context, roomID string) error {
	return r.dbSession.Query(deleteRoom, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
}

func (r RoomRepository) GetRoom(ctx context.Context, roomID string) (*domain.ChatRoom, error) {
	var room domain.ChatRoom
	studentMap := make(map[string]bool)
	var allStudents []domain.Student

	err := r.dbSession.Query(getRoom, roomID).WithContext(ctx).Consistency(gocql.One).Scan(&room.RoomID, &room.Admin.ID, &room.Class, &room.Deleted, &room.MaxParticipants, &room.Name, &studentMap)
	if err != nil {
		return nil, err
	}

	for userID, isPending := range studentMap {
		var student domain.Student
		student.ID = userID
		student.IsPending = isPending
		allStudents = append(allStudents, student)
	}
	room.Students = allStudents

	return &room, err
}

func (r RoomRepository) RemoveParticipantFromRoom(ctx context.Context, userID string, roomID string) error {
	return r.dbSession.Query(removeParticipantFromRoom, userID, roomID).WithContext(ctx).Consistency(gocql.One).Exec()
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
	retrievedChatRooms := make([]domain.ChatRoom, 0)
	var scanner cassandra.ScannerInterface
	studentMap := make(map[string]bool)
	scanner = r.dbSession.Query(getChatRoomsByClass, className).WithContext(ctx).Consistency(gocql.One).Iter().Scanner()

	for scanner.Next() {
		var room domain.ChatRoom
		err := scanner.Scan(&room.RoomID, &room.Admin.ID, &room.Class, &room.Deleted, &room.MaxParticipants, &room.Name, &studentMap)

		if err != nil {
			return nil, err
		}

		for userID, isPending := range studentMap {
			var student domain.Student
			student.ID = userID
			student.IsPending = isPending
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

	studentMap := make(map[string]bool)
	for _, student := range room.Students {
		studentMap[student.ID] = false
	}
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: saveRoom,
		Args: []interface{}{room.RoomID, room.Name, room.Admin.ID, studentMap, room.Class, room.MaxParticipants},
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
		Stmt: updateParticipantPendingState,
		Args: []interface{}{userID, false, roomID},
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
		Args: []interface{}{userID, roomID},
	})
	batch.AddBatchEntry(&gocql.BatchEntry{
		Stmt: removeRoomForParticipant,
		Args: []interface{}{[1]string{roomID}, userID},
	})
	return r.dbSession.ExecuteBatch(batch)
}
