package usecase

import (
	"bytes"
	"chat/domain"
	"chat/utils"
	"chat/utils/errors"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
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

func (u *messageUseCase) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time, userID string) (*domain.Message, error) {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	existingMessage, err := u.messageRepository.GetMessage(ctx, roomID, timeStamp)
	if err != nil {
		return nil, errors.NewNotFoundError("Message does not exist")
	}

	if userID != existingMessage.FromStudentID {
		return nil, errors.NewUnauthorizedError("Users can only delete their own messages")
	}

	err = u.messageRepository.DeleteMessage(c, roomID, timeStamp)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	return existingMessage, nil
}

func (u *messageUseCase) JoinRequest(ctx context.Context, roomID string, userID string, timeStamp time.Time) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	student, err := u.studentRepository.GetStudent(c, userID)
	if err != nil {
		return errors.NewNotFoundError(fmt.Sprintf("Student does not exist: %s", err.Error()))
	}

	room, err := u.roomRepository.GetRoom(c, roomID)
	if err != nil {
		return errors.NewConflictError(fmt.Sprintf("Room with ID %s does not exist: %s", roomID, err.Error()))
	}

	for _, participant := range room.Students {
		if participant.ID == userID {
			return errors.NewConflictError(fmt.Sprintf("User %s is already in room", userID))
		}
	}

	err = u.roomRepository.UpdateParticipantPendingState(c, roomID, userID, true)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Unable to update pending state: %s", err.Error()))
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

	err = u.roomRepository.RemoveParticipantFromRoom(c, userID, roomID)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Unable to remove user from room: %s", err.Error()))
	}

	emailBody, err := createEmailBody(student, roomID)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("Unable to create email: %s", err.Error()))
	}
	return u.mailer.SendSimpleMail(student.Email, emailBody)
}

func createEmailBody(student *domain.Student, team string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("Unable to retrieve current working directory\n %s", err))
	}

	var pathToFile string
	if strings.Contains(cwd, "bin") {
		pathToFile = path.Join(cwd, "rejection_template.html")
	} else {
		pathToFile = path.Join(cwd, "static", "rejection_template.html")
	}
	t, err := template.ParseFiles(pathToFile)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("Unable to find the file\n %s", err))
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
