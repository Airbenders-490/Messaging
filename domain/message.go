package domain

import (
	"context"
	"github.com/gocql/gocql"
	"time"
)

// Message struct
type Message struct {
	roomID gocql.UUID
	SentTimestamp time.Time
	FromStudentID string
	MessageBody string
}

// MessageRepository interface defines the functions all chatRepositories should have
type MessageRepository interface {
	SaveMessage(ctx context.Context, message *Message) error
	EditMessage(ctx context.Context, message *Message) error
	GetMessages(ctx context.Context, roomID string, timeStamp time.Time) ([]Message, error)
	DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error
}

// MessageUseCase defines the functionality messages encapsulate
type MessageUseCase interface {
	SendMessage(ctx context.Context, message *Message) error
	EditMessage(ctx context.Context, userID string, message *Message) error
	GetMessages(ctx context.Context) ([]Message, error)
	DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error
}