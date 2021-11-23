package repository

import (
	"chat/domain"
	"github.com/gocql/gocql"
)

type chatRepository struct {
	dbSession *gocql.Session
}

// NewChatRepository is the constructor
func NewChatRepository(session *gocql.Session) domain.ChatRepository {
	return &chatRepository{
		dbSession: session,
	}
}

