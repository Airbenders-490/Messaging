package usecase

import (
	"chat/domain"
	"chat/domain/mocks"
	"context"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const messageType = "*domain.Message"

func TestSaveMessage(t *testing.T) {
	t.Parallel()
	mockMessageRepository := new(mocks.MessageRepository)

	var mockMessage domain.Message
	faker.FakeData(&mockMessage)
	u := NewMessageUseCase(time.Second*2, mockMessageRepository, nil)

	t.Run("success", func(t *testing.T) {
		mockMessageRepository.
			On("SaveMessage", mock.Anything, mock.AnythingOfType(messageType)).
			Return(nil).Once()

		err := u.SaveMessage(context.TODO(), &mockMessage)

		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockMessageRepository.
			On("SaveMessage", mock.Anything, mock.AnythingOfType(messageType)).
			Return(errors.New("error")).Once()

		err := u.SaveMessage(context.TODO(), &mockMessage)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})
}

func TestEditMessage(t *testing.T) {
	t.Parallel()
	mockMessageRepository := new(mocks.MessageRepository)

	var mockMessage domain.Message
	faker.FakeData(&mockMessage)
	u := NewMessageUseCase(time.Second*2, mockMessageRepository, nil)

	t.Run("success", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()
		mockMessageRepository.
			On("EditMessage", mock.Anything, mock.AnythingOfType(messageType)).
			Return(nil).Once()

		editedMsg, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, "edited message")

		assert.NoError(t, err)

		assert.Equal(t, mockMessage.MessageBody, editedMsg.MessageBody)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("message is the same", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()

		editedMsg, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, mockMessage.MessageBody)

		assert.NoError(t, err)

		assert.EqualValues(t, mockMessage, *editedMsg)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("message does not exist", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(nil, errors.New("error")).Once()

		_, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, mockMessage.MessageBody)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("unauthorized user", func(t *testing.T) {
		var msg domain.Message
		faker.FakeData(&msg)
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&msg, nil).Once()

		_, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, mockMessage.MessageBody)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("message is empty", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()

		mockMessageRepository.
			On("DeleteMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(nil).Once()

		_, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, "")

		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("error editing message", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()
		mockMessageRepository.
			On("EditMessage", mock.Anything, mock.AnythingOfType(messageType)).
			Return(errors.New("error")).Once()

		_, err := u.EditMessage(context.TODO(), mockMessage.RoomID, mockMessage.FromStudentID,
			mockMessage.SentTimestamp, "editedMessage")

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})
}

func TestGetMessages(t *testing.T) {
	t.Parallel()
	mockMessageRepository := new(mocks.MessageRepository)
	var mockMessage []domain.Message

	faker.FakeData(&mockMessage)
	u := NewMessageUseCase(time.Second*2, mockMessageRepository, nil)

	t.Run("success", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessages", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("int")).
			Return(mockMessage[:1], nil).Once()

		retrievedMsgs, err := u.GetMessages(context.TODO(), mockMessage[0].RoomID, mockMessage[0].SentTimestamp, 1)

		assert.NotNil(t, retrievedMsgs)
		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)

	})

	t.Run("error", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessages", mock.Anything, mock.AnythingOfType("string"), mock.Anything, 5).
			Return(nil, errors.New("error")).Once()

		retrievedMsgs, err := u.GetMessages(context.TODO(), mockMessage[0].RoomID, mockMessage[0].SentTimestamp, 5)

		assert.Nil(t, retrievedMsgs)
		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

}

func TestDeleteMessage(t *testing.T) {
	t.Parallel()
	mockMessageRepository := new(mocks.MessageRepository)
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)
	u := NewMessageUseCase(time.Second*2, mockMessageRepository, nil)

	t.Run("success", func(t *testing.T) {
		mockMessageRepository.
			On("DeleteMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(nil).Once()

		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()

		err := u.DeleteMessage(context.TODO(), mockMessage.RoomID, mockMessage.SentTimestamp, mockMessage.FromStudentID)

		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("error: message does not exist", func(t *testing.T) {
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(nil, errors.New("error")).Once()

		err := u.DeleteMessage(context.TODO(), mockMessage.RoomID, mockMessage.SentTimestamp, mockMessage.FromStudentID)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("error unauthorized deletion", func(t *testing.T) {
		var msg domain.Message
		faker.FakeData(msg)
		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&msg, nil).Once()

		err := u.DeleteMessage(context.TODO(), mockMessage.RoomID, mockMessage.SentTimestamp, mockMessage.FromStudentID)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)

	})

	t.Run("error unable to delete", func(t *testing.T) {
		mockMessageRepository.
			On("DeleteMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(errors.New("error")).Once()

		mockMessageRepository.
			On("GetMessage", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
			Return(&mockMessage, nil).Once()

		err := u.DeleteMessage(context.TODO(), mockMessage.RoomID, mockMessage.SentTimestamp, mockMessage.FromStudentID)

		assert.Error(t, err)

		mockMessageRepository.AssertExpectations(t)
	})
}
