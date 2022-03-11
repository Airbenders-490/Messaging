// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	domain "chat/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MessageUseCase is an autogenerated mock type for the MessageUseCase type
type MessageUseCase struct {
	mock.Mock
}

// DeleteMessage provides a mock function with given fields: ctx, roomID, timeStamp, userID
func (_m *MessageUseCase) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time, userID string) error {
	ret := _m.Called(ctx, roomID, timeStamp, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, string) error); ok {
		r0 = rf(ctx, roomID, timeStamp, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditMessage provides a mock function with given fields: ctx, roomID, userID, timeStamp, message
func (_m *MessageUseCase) EditMessage(ctx context.Context, roomID string, userID string, timeStamp time.Time, message string) (*domain.Message, error) {
	ret := _m.Called(ctx, roomID, userID, timeStamp, message)

	var r0 *domain.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time, string) *domain.Message); ok {
		r0 = rf(ctx, roomID, userID, timeStamp, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, time.Time, string) error); ok {
		r1 = rf(ctx, roomID, userID, timeStamp, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessages provides a mock function with given fields: ctx, roomID, timeStamp, limit
func (_m *MessageUseCase) GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]domain.Message, error) {
	ret := _m.Called(ctx, roomID, timeStamp, limit)

	var r0 []domain.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, int) []domain.Message); ok {
		r0 = rf(ctx, roomID, timeStamp, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, int) error); ok {
		r1 = rf(ctx, roomID, timeStamp, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsAuthorized provides a mock function with given fields: ctx, userID, roomID
func (_m *MessageUseCase) IsAuthorized(ctx context.Context, userID string, roomID string) bool {
	ret := _m.Called(ctx, userID, roomID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, userID, roomID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// JoinRequest provides a mock function with given fields: ctx, roomID, userID, timeStamp
func (_m *MessageUseCase) JoinRequest(ctx context.Context, roomID string, userID string, timeStamp time.Time) error {
	ret := _m.Called(ctx, roomID, userID, timeStamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) error); ok {
		r0 = rf(ctx, roomID, userID, timeStamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMessage provides a mock function with given fields: ctx, message
func (_m *MessageUseCase) SaveMessage(ctx context.Context, message *domain.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendRejection provides a mock function with given fields: ctx, roomID, userID, loggedID
func (_m *MessageUseCase) SendRejection(ctx context.Context, roomID string, userID string, loggedID string) error {
	ret := _m.Called(ctx, roomID, userID, loggedID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, roomID, userID, loggedID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
