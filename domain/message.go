package domain

import (
	"context"
	"time"
)

// Message struct
type Message struct {
	RoomID        string
	SentTimestamp time.Time
	FromStudentID string
	MessageBody   string
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
	SaveMessage(ctx context.Context, message *Message) error
	EditMessage(ctx context.Context, userID string, message *Message) error
	GetMessages(ctx context.Context, roomID string, timeStamp time.Time) ([]Message, error)
	DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error
	IsAuthorized(userID, roomID string) bool
}
