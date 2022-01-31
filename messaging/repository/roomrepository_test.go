package repository

import (
	"chat/domain"
	"chat/messaging/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
)

var batchMock = &mocks.BatchInterface{}
var sessionMock = &mocks.SessionInterface{}
var queryMock = &mocks.QueryInterface{}
var ctx = context.Background()
var rr = NewRoomRepository(sessionMock)
var room = &domain.ChatRoom{
	Students: []domain.Student{{ID: "userID1"}},
}

func resetFields() {
	batchMock = &mocks.BatchInterface{}
	sessionMock = &mocks.SessionInterface{}
	queryMock = &mocks.QueryInterface{}
	ctx = context.Background()
	rr = NewRoomRepository(sessionMock)
	room = &domain.ChatRoom{
		Students: []domain.Student{{ID: "userID1"}},
	}
}

func TestAddParticipantToRoomSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.AddParticipantToRoom(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestDeleteRoomSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.DeleteRoom(ctx, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomSuccessEmptyRoom(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := rr.GetRoom(ctx, mock.Anything)

	if err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Error"))

	_, err := rr.GetRoom(ctx, mock.Anything)

	if err == nil {
		t.Errorf("Actual no error, expected error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveParticipantFromRoomSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.RemoveParticipantFromRoom(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestSaveRoomSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.SaveRoom(ctx, room); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestAddRoomForParticipantSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.AddRoomForParticipant(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestAddRoomForParticipantsSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.AddRoomForParticipants(ctx, mock.Anything, []string{"userID1", "userID2"}); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	sessionMock.AssertNumberOfCalls(t, "Query", 2)
	resetFields()
}

func TestAddRoomForParticipantsFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(errors.New("Internal Error"))

	if err := rr.AddRoomForParticipants(ctx, mock.Anything, []string{"userID1"}); err == nil {
		t.Errorf("Actual no error, expected error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomsForSuccessNoMatches(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything).Return(nil)

	_, err := rr.GetRoomsFor(ctx, mock.Anything)

	if err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomsForFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything).Return(errors.New("Internal Error"))

	_, err := rr.GetRoomsFor(ctx, mock.Anything)

	if err == nil {
		t.Errorf("Actual no error, expected error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.RemoveRoomForParticipant(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantsSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.RemoveRoomForParticipants(ctx, mock.Anything, []domain.Student{{"userID1", "first1", "last1"}}); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantsFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(errors.New("Internal Error"))

	if err := rr.RemoveRoomForParticipants(ctx, mock.Anything, []domain.Student{{"userID1", "first1", "last1"}}); err == nil {
		t.Errorf("Actual no error, expected error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestSaveRoomAndAddRoomForAllParticipantsSuccess(t *testing.T) {
	sessionMock.On("NewBatch", mock.Anything).Return(batchMock)
	batchMock.On("WithContext", ctx).Return(batchMock)
	batchMock.On("AddBatchEntry", mock.Anything)
	sessionMock.On("ExecuteBatch", batchMock).Return(nil)

	if err := rr.SaveRoomAndAddRoomForAllParticipants(ctx, room); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantsAndDeleteRoomSuccess(t *testing.T) {
	sessionMock.On("NewBatch", mock.Anything).Return(batchMock)
	batchMock.On("WithContext", ctx).Return(batchMock)
	batchMock.On("AddBatchEntry", mock.Anything)
	sessionMock.On("ExecuteBatch", batchMock).Return(nil)

	if err := rr.RemoveRoomForParticipantsAndDeleteRoom(ctx, room); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestAddParticipantToRoomAndAddRoomForParticipantSuccess(t *testing.T) {
	sessionMock.On("NewBatch", mock.Anything).Return(batchMock)
	batchMock.On("WithContext", ctx).Return(batchMock)
	batchMock.On("AddBatchEntry", mock.Anything)
	sessionMock.On("ExecuteBatch", batchMock).Return(nil)

	if err := rr.AddParticipantToRoomAndAddRoomForParticipant(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveParticipantFromRoomAndRemoveRoomForParticipantSuccess(t *testing.T) {
	sessionMock.On("NewBatch", mock.Anything).Return(batchMock)
	batchMock.On("WithContext", ctx).Return(batchMock)
	batchMock.On("AddBatchEntry", mock.Anything)
	sessionMock.On("ExecuteBatch", batchMock).Return(nil)

	if err := rr.RemoveParticipantFromRoomAndRemoveRoomForParticipant(ctx, mock.Anything, mock.Anything); err != nil {
		t.Errorf("Actual error, expected no error")
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}
