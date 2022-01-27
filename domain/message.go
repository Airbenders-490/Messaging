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
	GetMessage(ctx context.Context, roomID string, timeStamp time.Time) (*Message, error)
	GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]Message, error)
	DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error
}

// MessageUseCase defines the functionality messages encapsulate
type MessageUseCase interface {
	SaveMessage(ctx context.Context, message *Message) error
	EditMessage(ctx context.Context, roomID string, userID string, timeStamp time.Time, message string) (*Message, error)
	GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]Message, error)
	DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error
	IsAuthorized(ctx context.Context, userID, roomID string) bool
}
