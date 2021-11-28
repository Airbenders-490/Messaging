package repository

import (
	"chat/domain"
	"context"
	"github.com/gocql/gocql"
	"time"
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

func (m messageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	panic("implement me")
}

func (m messageRepository) EditMessage(ctx context.Context, message *domain.Message) error {
	panic("implement me")
}

func (m messageRepository) GetMessages(ctx context.Context, roomID string, timeStamp time.Time) ([]domain.Message, error) {
	panic("implement me")
}

func (m messageRepository) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error {
	panic("implement me")
}


