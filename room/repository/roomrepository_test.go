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
var mockIter = &mocks.IterInterface{}
var scannerMock = &mocks.ScannerInterface{}
var ctx = context.Background()
var rr = NewRoomRepository(sessionMock)
var room = &domain.ChatRoom{
	Students: []domain.Student{{ID: "userID1"}},
}

const errorMessage = "Actual error, expected no error"
const internalErrorMessage = "Internal Error"
const errorMessage2 = "Actual no error, expected error"

func resetFields() {
	batchMock = &mocks.BatchInterface{}
	sessionMock = &mocks.SessionInterface{}
	queryMock = &mocks.QueryInterface{}
	mockIter = &mocks.IterInterface{}
	ctx = context.Background()
	rr = NewRoomRepository(sessionMock)
	room = &domain.ChatRoom{
		Students: []domain.Student{{ID: "userID1"}},
	}
}

func TestDeleteRoomSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.DeleteRoom(ctx, mock.Anything); err != nil {
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomSuccessEmptyRoom(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := rr.GetRoom(ctx, mock.Anything)

	if err != nil {
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New(internalErrorMessage))

	_, err := rr.GetRoom(ctx, mock.Anything)

	if err == nil {
		t.Errorf(errorMessage2)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	sessionMock.AssertNumberOfCalls(t, "Query", 2)
	resetFields()
}

func TestAddRoomForParticipantsFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(errors.New(internalErrorMessage))

	if err := rr.AddRoomForParticipants(ctx, mock.Anything, []string{"userID1"}); err == nil {
		t.Errorf(errorMessage2)
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
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetRoomsForFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Scan", mock.Anything, mock.Anything).Return(errors.New(internalErrorMessage))

	_, err := rr.GetRoomsFor(ctx, mock.Anything)

	if err == nil {
		t.Errorf(errorMessage2)
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
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantsSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.RemoveRoomForParticipants(ctx, mock.Anything, []domain.Student{{"userID1", "", "", ""}}); err != nil {
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestRemoveRoomForParticipantsFail(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(errors.New(internalErrorMessage))

	if err := rr.RemoveRoomForParticipants(ctx, mock.Anything, []domain.Student{{"userID1", "", "", ""}}); err == nil {
		t.Errorf(errorMessage2)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
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
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetChatRoomsByClassSuccess(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Iter").Return(mockIter)
	mockIter.On("Scanner").Return(scannerMock)
	scannerMock.On("Next").Return(true).Once()
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	scannerMock.On("Next").Return(false).Once()
	scannerMock.On("Err").Return(nil).Once()

	if _, err := rr.GetChatRoomsByClass(ctx, mock.Anything); err != nil {
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetChatRoomsByClassFailScan(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Iter").Return(mockIter)
	mockIter.On("Scanner").Return(scannerMock)
	scannerMock.On("Next").Return(true).Once()

	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New(internalErrorMessage)).Once()

	if _, err := rr.GetChatRoomsByClass(ctx, mock.Anything); err == nil {
		t.Errorf(errorMessage2)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestGetChatRoomsByClassFailCloseScan(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Iter").Return(mockIter)
	mockIter.On("Scanner").Return(scannerMock)
	scannerMock.On("Next").Return(true).Once()
	scannerMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	scannerMock.On("Next").Return(false).Once()
	scannerMock.On("Err").Return(errors.New(internalErrorMessage))

	if _, err := rr.GetChatRoomsByClass(ctx, mock.Anything); err == nil {
		t.Errorf(errorMessage2)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}

func TestUpdateParticipantPendingState(t *testing.T) {
	sessionMock.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("bool"), mock.AnythingOfType("string")).Return(queryMock)
	queryMock.On("WithContext", ctx).Return(queryMock)
	queryMock.On("Consistency", mock.Anything).Return(queryMock)
	queryMock.On("Exec").Return(nil)

	if err := rr.UpdateParticipantPendingState(ctx, "", "", false); err != nil {
		t.Errorf(errorMessage)
	}
	sessionMock.AssertExpectations(t)
	resetFields()
}
