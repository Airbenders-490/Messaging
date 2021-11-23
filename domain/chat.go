package domain

import "time"

// Chat struct
type Chat struct {
	MessageID string
	FromStudentID string
	ToStudentID string
	MessageBody string
	TeamID string
	SentTimestamp time.Time
}

// ChatRepository interface defines the functions all chatRepositories should have
type ChatRepository interface {

}