package repository

import (
	"chat/domain"
	"chat/messaging/repository/cassandra"
	"chat/utils/errors"
	"context"
	"time"
)

const (
	insertMessage = `INSERT INTO chat.messages (room_id, from_student_id, message_body, sent_timestamp) VALUES (?, ?, ?, ?)`
	editMessage   = `UPDATE chat.messages SET message_body=? WHERE room_id=? AND sent_timestamp=? IF EXISTS;`
	getMessage    = `SELECT * FROM chat.messages where room_id=? AND sent_timestamp =?`
	getMessages   = `SELECT * FROM chat.messages WHERE room_id=? AND sent_timestamp <? limit ?`
	deleteMessage = `DELETE FROM chat.messages WHERE room_id=? AND sent_timestamp=? IF EXISTS`
)

type MessageRepository struct {
	dbSession cassandra.SessionInterface
}

// NewChatRepository is the constructor
func NewChatRepository(session cassandra.SessionInterface) *MessageRepository {
	return &MessageRepository{
		dbSession: session,
	}
}

func (m *MessageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	return m.dbSession.Query(insertMessage, message.RoomID, message.FromStudentID, message.MessageBody, message.SentTimestamp).WithContext(ctx).Exec()
}

func (m *MessageRepository) EditMessage(ctx context.Context, message *domain.Message) error {
	var editedMsg domain.Message
	applied, err := m.dbSession.Query(editMessage, message.MessageBody, message.RoomID, message.SentTimestamp).WithContext(ctx).
		ScanCAS(&editedMsg.RoomID, &editedMsg.SentTimestamp, &editedMsg.FromStudentID, &editedMsg.MessageBody)

	if err != nil {
		return err
	}

	if !applied {
		return errors.NewInternalServerError("Unable to save the message")
	}

	return nil
}

func (m *MessageRepository) GetMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	var retrievedMsg domain.Message

	err := m.dbSession.Query(getMessage, message.RoomID, message.SentTimestamp).WithContext(ctx).
		Scan(&retrievedMsg.RoomID, &retrievedMsg.SentTimestamp, &retrievedMsg.FromStudentID, &retrievedMsg.MessageBody)
	if err != nil {
		return nil, err
	}

	return &retrievedMsg, err
}

func (m *MessageRepository) GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]domain.Message, error) {
	var retrievedMessages []domain.Message
	var scanner cassandra.ScannerInterface

	scanner = m.dbSession.Query(getMessages, roomID, timeStamp, limit).WithContext(ctx).Iter().Scanner()

	for scanner.Next() {
		var msg domain.Message
		err := scanner.Scan(&msg.RoomID, &msg.SentTimestamp, &msg.FromStudentID, &msg.MessageBody)

		if err != nil {
			return nil, err
		}
		retrievedMessages = append(retrievedMessages, msg)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return retrievedMessages, nil
}

func (m *MessageRepository) DeleteMessage(ctx context.Context, message *domain.Message) error {
	return m.dbSession.Query(deleteMessage, message.RoomID, message.SentTimestamp).WithContext(ctx).Exec()
}
