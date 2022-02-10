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

	return
}

func (u *messageUseCase) SaveMessage(ctx context.Context, message *domain.Message) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.messageRepository.SaveMessage(c, message)
}

func (u *messageUseCase) EditMessage(ctx context.Context, roomID string, userID string, timeStamp time.Time, message string) (*domain.Message, error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	existingMessage, err := u.messageRepository.GetMessage(ctx, roomID, timeStamp)
	if err != nil {
		return nil, errors.NewNotFoundError("Message does not exist")
	}

	if userID != existingMessage.FromStudentID {
		return nil, errors.NewUnauthorizedError("Users can only edit their own messages")
	}

	if message == existingMessage.MessageBody {
		return existingMessage, nil
	}

	if message == "" {
		return nil, u.messageRepository.DeleteMessage(c, roomID, timeStamp)
	}

	existingMessage.MessageBody = message
	err = u.messageRepository.EditMessage(c, existingMessage)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	return existingMessage, nil
}

func (u *messageUseCase) GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]domain.Message, error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	retrievedMessages, err := u.messageRepository.GetMessages(c, roomID, timeStamp, limit)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return retrievedMessages, nil
}

func (u *messageUseCase) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time, userID string) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	existingMessage, err := u.messageRepository.GetMessage(ctx, roomID, timeStamp)
	if err != nil {
		return errors.NewNotFoundError("Message does not exist")
	}

	if userID != existingMessage.FromStudentID {
		return errors.NewUnauthorizedError("Users can only delete their own messages")
	}

	return u.messageRepository.DeleteMessage(c, roomID, timeStamp)
}
