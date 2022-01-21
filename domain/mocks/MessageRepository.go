// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	domain "chat/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

// DeleteMessage provides a mock function with given fields: ctx, roomID, timeStamp
func (_m *MessageRepository) DeleteMessage(ctx context.Context, roomID string, timeStamp time.Time) error {
	ret := _m.Called(ctx, roomID, timeStamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, roomID, timeStamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EditMessage provides a mock function with given fields: ctx, message
func (_m *MessageRepository) EditMessage(ctx context.Context, message *domain.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMessage provides a mock function with given fields: ctx, roomID, timeStamp
func (_m *MessageRepository) GetMessage(ctx context.Context, roomID string, timeStamp time.Time) (*domain.Message, error) {
	ret := _m.Called(ctx, roomID, timeStamp)

	var r0 *domain.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) *domain.Message); ok {
		r0 = rf(ctx, roomID, timeStamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) error); ok {
		r1 = rf(ctx, roomID, timeStamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessages provides a mock function with given fields: ctx, roomID, timeStamp, limit
func (_m *MessageRepository) GetMessages(ctx context.Context, roomID string, timeStamp time.Time, limit int) ([]domain.Message, error) {
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

// SaveMessage provides a mock function with given fields: ctx, message
func (_m *MessageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}