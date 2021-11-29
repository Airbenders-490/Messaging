package repository

import (
	"chat/domain"
	"context"
	"github.com/gocql/gocql"
	"time"
)

const (
	insertMessage = `INSERT INTO messages (room_id, from_student_id, message_body, sent_timestamp) VALUES (?, ?, ?, ?)`
)

type messageRepository struct {
	dbSession *gocql.Session
}

// NewChatRepository is the constructor
func NewChatRepository(session *gocql.Session) domain.MessageRepository {
	return &messageRepository{
		dbSession: session,
	}
}

func (m *messageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	err := m.dbSession.Query(insertMessage, message.RoomID, message.FromStudentID, message.MessageBody, message.SentTimestamp).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *messageRepository) EditMessage(ctx context.Context, message *domain.Message) error {
	panic("implement me")
}

func (m *messageRepository) GetMessages(ctx context.Context, roomID string, timeStamp time.Time) ([]domain.Message, error) {
	panic("implement me")
}

func (m *messageRepository) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error {
	panic("implement me")
}
