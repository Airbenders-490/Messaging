package repository

import (
	"chat/domain"
	"chat/messaging/repository/mocks"
	"context"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/mock"
	"testing"
)

var sr = NewStudentRepository(sessionMock)
var mockStudent domain.Student

func resetStudentRepoFields() {
	sessionMock = &mocks.SessionInterface{}
	queryMock = &mocks.QueryInterface{}
	faker.FakeData(&mockStudent)
	ctx = context.Background()
	sr = NewStudentRepository(sessionMock)
}

func TestSaveStudentSuccess(t *testing.T) {
	resetStudentRepoFields()

	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := sr.SaveStudent(ctx, &mockStudent); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
}

func TestEditStudentSuccess(t *testing.T) {
	resetStudentRepoFields()

	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := sr.EditStudent(ctx, &mockStudent); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
}

func TestDeleteStudentSuccess(t *testing.T) {
	resetStudentRepoFields()

	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := sr.DeleteStudent(ctx, ""); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
}

func TestGetStudentSuccess(t *testing.T) {
	resetStudentRepoFields()

	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	if _, err := sr.GetStudent(ctx, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
}

func TestGetStudentFail(t *testing.T) {
	resetStudentRepoFields()

	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(""))

	if _, err := sr.GetStudent(ctx, mock.Anything); err == nil {
		t.Errorf("Actual no error, expected  error")
	}
	sessionMock.AssertExpectations(t)
}
