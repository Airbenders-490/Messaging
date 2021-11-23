package repository

import (
	"chat/domain"
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

type chatRepository struct {
	dbSession *gocql.Session
}

func NewChatRepository(session *gocql.Session) domain.ChatRepository {
	return &chatRepository{
		dbSession: session,
	}
}

const (
	selectByFromAndToID = `SELECT * FROM chat.messages WHERE from_student_id=? AND to_student_id=?;`
)

func(chatRepository *chatRepository) GetByFromAndToID(fromStudentID string, toStudentID string) (*domain.Chat, error) {

	var chat domain.Chat

	if err := chatRepository.dbSession.Query(selectByFromAndToID, fromStudentID, toStudentID).Consistency(gocql.One).Scan(&chat.FromStudentID, &chat.FromStudentID, &chat.SentTimestamp, &chat.MessageBody, &chat.MessageID, &chat.TeamID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message entry: ", chat.MessageID, chat.FromStudentID, chat.ToStudentID, chat.MessageBody, chat.TeamID, chat.SentTimestamp)

	return &chat, nil
}