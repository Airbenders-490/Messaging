package usecase

import (
	"chat/domain"
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

func (u *messageUseCase) IsAuthorized(userID, roomID string) (authorized bool) {
	_, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	studentChatRooms, err := u.roomRepository.GetRoomsFor(userID)
	if err != nil {
		return false
	}

	// todo: this is temporary. Must be removed!
	authorized = true
	return
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

func (u *messageUseCase) EditMessage(ctx context.Context, userID string, message *domain.Message) error {
	panic("")
}

func (u *messageUseCase) GetMessages(ctx context.Context, roomID string, timeStamp time.Time) ([]domain.Message, error) {
	panic("")
}

func (u *messageUseCase) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error {
	panic("")
}
