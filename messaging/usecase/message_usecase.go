package usecase

import (
	"bytes"
	"chat/domain"
	"chat/utils"
	"chat/utils/errors"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"
)

type messageUseCase struct {
	timeout           time.Duration
	mailer            utils.Mailer
	messageRepository domain.MessageRepository
	roomRepository    domain.RoomRepository
	studentRepository domain.StudentRepository
}

// NewMessageUseCase instantiates a
func NewMessageUseCase(
	t time.Duration,
	mr domain.MessageRepository,
	rr domain.RoomRepository,
	sr domain.StudentRepository,
	mailer utils.Mailer) domain.MessageUseCase {
	return &messageUseCase{timeout: t, messageRepository: mr, roomRepository: rr, studentRepository: sr, mailer: mailer}
}

func (u *messageUseCase) IsAuthorized(ctx context.Context, userID, roomID string) (authorized bool) {
	_, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	return true
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

func (u *messageUseCase) JoinRequest(ctx context.Context, roomID string, userID string, timeStamp time.Time) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	student, err := u.studentRepository.GetStudent(c, userID)
	if err != nil {
		return errors.NewNotFoundError("Student does not exist")
	}

	_, err = u.roomRepository.GetRoom(c, roomID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does not exist", roomID))
	}

	m := domain.Message{
		RoomID:        roomID,
		SentTimestamp: timeStamp,
		FromStudentID: userID,
		MessageBody:   fmt.Sprintf("%s %s has requested to join your group.", student.FirstName, student.LastName)}

	return u.messageRepository.SaveMessage(c, &m)
}

func (u *messageUseCase) SendRejection(ctx context.Context, roomID string, userID string, loggedID string) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	room, err := u.roomRepository.GetRoom(c, roomID)
	if err != nil {
		return errors.NewNotFoundError("Room does not exist")
	}

	if loggedID != room.Admin.ID {
		return errors.NewUnauthorizedError("You are not authorized to reject because you are not admin.")
	}

	student, err := u.studentRepository.GetStudent(c, userID)
	if err != nil {
		return errors.NewNotFoundError("User does not exist")
	}

	emailBody, err := createEmailBody(student, roomID)
	if err != nil {
		return err
	}
	return u.mailer.SendSimpleMail(student.Email, emailBody)
}

func createEmailBody(student *domain.Student, team string) ([]byte, error) {
	t, err := template.ParseFiles("./static/rejection_template.html")
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("Unable to find the file %s", err))
	}

	var body bytes.Buffer

	message := fmt.Sprintf("From: %s\r\n", os.Getenv("EMAIL_FROM"))
	message += fmt.Sprintf("To: %s\r\n", student.Email)
	message += "Subject: Team Request\r\n"
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message += "\r\n"

	body.Write([]byte(message))
	t.Execute(&body, struct {
		Name string
		Team string
	}{
		Name: student.FirstName,
		Team: team,
	})

	return body.Bytes(), nil
}
