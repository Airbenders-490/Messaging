package domain

import "time"

// Chat struct
type Message struct {
	roomID string
	SentTimestamp time.Time
	FromStudentID string
	MessageBody string
}

// ChatRepository interface defines the functions all chatRepositories should have
type ChatRepository interface {

}