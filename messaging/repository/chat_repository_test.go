package repository

import (
	"chat/domain"
	"chat/messaging/repository/mocks"
	"context"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var query = &mocks.QueryInterface{}
var session = &mocks.SessionInterface{}
var cr = NewChatRepository(session)
var iter = &mocks.IterInterface{}
var scannerMock = &mocks.ScannerInterface{}

func reset() {
	query = &mocks.QueryInterface{}
	session = &mocks.SessionInterface{}
	iter = &mocks.IterInterface{}
	scannerMock = &mocks.ScannerInterface{}
	cr = NewChatRepository(session)
}

const errorMessage = "Actual error, expected no error"
const internalErrorMessage = "Internal Error"
const errorMessage2 = "Actual no error, expected error"

func TestSaveMessageSuccess(t *testing.T){

	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Exec").
		Return(nil)

	err := cr.SaveMessage(context.Background(), &mockMessage)

	assert.NoError(t, err)

	session.AssertExpectations(t)
	//reset()
}

func TestSaveMessageError(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Exec").
		Return(errors.New(internalErrorMessage))

	err := cr.SaveMessage(context.Background(), &mockMessage)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestEditMessageSuccess(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("ScanCAS", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(true, nil).
		Once()

	err := cr.EditMessage(context.Background(), &mockMessage)

	assert.NoError(t, err)

	session.AssertExpectations(t)
}

func TestEditMessageFail(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("ScanCAS", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(false, errors.New("error")).
		Once()

	err := cr.EditMessage(context.Background(), &mockMessage)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestEditMessageUnableToUpdate(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("ScanCAS", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(false, nil).
		Once()

	err := cr.EditMessage(context.Background(), &mockMessage)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestGetMessageSuccess(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	_, err := cr.GetMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)

	assert.NoError(t, err)

	session.AssertExpectations(t)
}

func TestGetMessageError(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error"))

	_, err := cr.GetMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestGetMessagesSuccess(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Iter").
		Return(iter)
	iter.On("Scanner").
		Return(scannerMock)

	scannerMock.On("Next").Return(true).Once()
	scannerMock.On("Next").Return(false)
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	scannerMock.On("Err").Return(nil)

	_, err := cr.GetMessages(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp, 2)

	assert.NoError(t, err)

	session.AssertExpectations(t)
}

func TestGetMessagesScanError(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Iter").
		Return(iter)
	iter.On("Scanner").
		Return(scannerMock)

	scannerMock.On("Next").Return(true).Once()
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New(internalErrorMessage))

	_, err := cr.GetMessages(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp, 2)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestGetMessagesError(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Iter").
		Return(iter)
	iter.On("Scanner").
		Return(scannerMock)

	scannerMock.On("Next").Return(false).Once()
	scannerMock.On("Err").Return(errors.New("error"))

	_, err := cr.GetMessages(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp, 2)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

func TestDeleteMessageSuccess(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Exec").Return(nil)

	err := cr.DeleteMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)

	assert.NoError(t, err)

	session.AssertExpectations(t)
}

func TestDeleteMessageError(t *testing.T) {
	reset()
	var mockMessage domain.Message
	faker.FakeData(&mockMessage)

	session.On("Query", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(query)
	query.On("WithContext", mock.Anything).
		Return(query)
	query.On("Exec").Return(errors.New(internalErrorMessage))

	err := cr.DeleteMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)

	assert.Error(t, err)

	session.AssertExpectations(t)
}

//func TestDeleteMessage(t *testing.T){
//	t.Parallel()
//	faker.FakeData(&mockMessage)
//	reset()
//	t.Run("succes", func(t *testing.T) {
//		session.On("Query",mock.AnythingOfType("string"), mock.Anything, mock.Anything).
//			Return(query)
//		query.On("WithContext", mock.Anything).
//			Return(query)
//		query.On("Exec").Return(nil)
//
//		err := cr.DeleteMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)
//
//		assert.NoError(t, err)
//
//		session.AssertExpectations(t)
//		reset()
//	})
//
//	t.Run("succes", func(t *testing.T) {
//		session.On("Query",mock.AnythingOfType("string"), mock.Anything, mock.Anything).
//			Return(query)
//		query.On("WithContext", mock.Anything).
//			Return(query)
//		query.On("Exec").Return(errors.New(internalErrorMessage))
//
//		err := cr.DeleteMessage(context.Background(), mockMessage.RoomID, mockMessage.SentTimestamp)
//
//		assert.Error(t, err)
//
//		session.AssertExpectations(t)
//		reset()
//	})
//}
