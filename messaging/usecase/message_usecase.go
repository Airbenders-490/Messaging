package usecase

import (
	"chat/domain"
	"chat/utils/errors"
	"context"
	"time"
)

type messageUseCase struct {
	timeout           time.Duration
	messageRepository domain.MessageRepository
	roomRepository    domain.RoomRepository
}

// NewMessageUseCase instantiates a
func NewMessageUseCase(t time.Duration, mr domain.MessageRepository, rr domain.RoomRepository) domain.MessageUseCase {
	return &messageUseCase{timeout: t, messageRepository: mr, roomRepository: rr}
}

func (u *messageUseCase) IsAuthorized(ctx context.Context, userID, roomID string) (authorized bool) {
	_, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	studentChatRooms, err := u.roomRepository.GetRoomsFor(ctx, userID)
	if err != nil {
		return false
	}

	for _, room := range studentChatRooms.Rooms {
		if roomID == room.RoomID {
			authorized = true
			break
		}
	}

	return authorized
}

func (u *messageUseCase) SaveMessage(ctx context.Context, message *domain.Message) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	err := u.messageRepository.SaveMessage(c, message)
	if err != nil {
		return err
	}

	return nil
}

func (u *messageUseCase) EditMessage(ctx context.Context, roomID string, userID string, timeStamp time.Time, message string) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	//message.SentTimestamp, message.FromStudentID, message.RoomID = timeStamp, userID, roomID
	//get existing message at timestamp
	existingMessage, err := u.messageRepository.GetMessage(ctx, roomID, timeStamp)
	if err != nil {
		return err
	}

	//if userID != existingMessage.FromStudentID{
	//	return errors.NewUnauthorizedError("Users can only edit their own messages")
	//}

	if message == "" || message == existingMessage.MessageBody {
		return errors.NewBadRequestError("Message body is the same or empty")
	}

	existingMessage.MessageBody = message
	err = u.messageRepository.EditMessage(c, existingMessage)

	if err != nil {
		return err
	}

	return nil
}

func (u *messageUseCase) GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]domain.Message, error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if retrievedMessages, err := u.messageRepository.GetMessages(c, roomID, timeStamp, limit); err != nil{
		return nil, err
	}else{
		return retrievedMessages, nil
	}


}
func (u *messageUseCase) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if err := u.messageRepository.DeleteMessage(c, roomID, timeStamp); err != nil{
		return err
	}

	return nil
}
