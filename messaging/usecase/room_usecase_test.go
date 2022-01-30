package usecase

import (
	"chat/domain"
	"chat/domain/mocks"
	"context"
	"errors"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSaveRoom(t *testing.T) {
	mockRoomRepo := new(mocks.RoomRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockRoom domain.ChatRoom
	faker.FakeData(&mockRoom)

	t.Run("case success", func(t *testing.T) {
		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil)
		mockRoomRepo.On("SaveRoomAndAddRoomForAllParticipants",mock.Anything,mock.Anything).
			Return(nil).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.SaveRoom(context.TODO(), &mockRoom)
		assert.NoError(t, err)
		mockRoomRepo.AssertExpectations(t)

	})

	t.Run("case room exists", func(t *testing.T) {
		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(&mockRoom, nil).
			Once()

		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.SaveRoom(context.TODO(), &mockRoom)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})

	t.Run("case student does not exist", func(t *testing.T) {
		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(&mockRoom, nil).
			Once()
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(nil,errors.New("error")).Once()

		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.SaveRoom(context.TODO(), &mockRoom)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
	t.Run("case error in repo", func(t *testing.T) {
		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil)
		mockRoomRepo.On("SaveRoomAndAddRoomForAllParticipants",mock.Anything,mock.Anything).
			Return(errors.New("error")).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.SaveRoom(context.TODO(), &mockRoom)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
}
func TestAddUserToRoom(t *testing.T) {
	mockRoomRepo := new(mocks.RoomRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockRoom domain.ChatRoom
	faker.FakeData(&mockRoom)

	t.Run("case success", func(t *testing.T) {
		mockRoomRepo.On("AddParticipantToRoomAndAddRoomForParticipant",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).
			Return(nil).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.AddUserToRoom(context.TODO(),mockRoom.RoomID,mockStudent.ID)
		assert.NoError(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
	t.Run("case error in repo", func(t *testing.T) {
		mockRoomRepo.On("AddParticipantToRoomAndAddRoomForParticipant",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).
			Return(errors.New("error")).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.AddUserToRoom(context.TODO(),mockRoom.RoomID,mockStudent.ID)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
}

func TestRemoveUserFromRoom(t *testing.T) {
	mockRoomRepo := new(mocks.RoomRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockRoom domain.ChatRoom
	faker.FakeData(&mockRoom)

	t.Run("case success", func(t *testing.T) {
		mockRoomRepo.On("RemoveParticipantFromRoomAndRemoveRoomForParticipant",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).
			Return(nil).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.RemoveUserFromRoom(context.TODO(),mockRoom.RoomID,mockStudent.ID)
		assert.NoError(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
	t.Run("case error in repo", func(t *testing.T) {
		mockRoomRepo.On("RemoveParticipantFromRoomAndRemoveRoomForParticipant",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).
			Return(errors.New("error")).Once()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		err:=u.RemoveUserFromRoom(context.TODO(),mockRoom.RoomID,mockStudent.ID)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
}
func TestGetChatRoomsFor(t *testing.T) {
	mockRoomRepo := new(mocks.RoomRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockRoom domain.ChatRoom
	faker.FakeData(&mockRoom)
	var mockStudentChatRoom domain.StudentChatRooms
	faker.FakeData(&mockStudentChatRoom)
	t.Run("case room in rooms for student does not exist", func(t *testing.T) {
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil).Once()

		mockRoomRepo.On("GetRoomsFor",mock.Anything, mock.Anything).
			Return(&mockStudentChatRoom, nil).Once()

		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(nil,errors.New("error")).Maybe()



		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		chatroom,err:=u.GetChatRoomsFor(context.TODO(),mockStudent.ID)
		assert.Error(t, err)
		assert.Nil(t, chatroom)

		mockRoomRepo.AssertExpectations(t)
	})
	t.Run("case student in room does not exist", func(t *testing.T) {
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
		Return(&mockStudent,nil).Once()

		mockRoomRepo.On("GetRoomsFor",mock.Anything, mock.Anything).
		Return(&mockStudentChatRoom, nil).Once()

		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
		Return(&mockRoom,nil)

		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
		Return(nil,errors.New("error"))

		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		chatroom,err:=u.GetChatRoomsFor(context.TODO(),mockStudent.ID)
		assert.Error(t, err)
		assert.Nil(t, chatroom)

		mockRoomRepo.AssertExpectations(t)
	})

	t.Run("case error in repo or chatroom does not exist", func(t *testing.T) {
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil).Once()

		mockRoomRepo.On("GetRoomsFor",mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).Once()

		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(&mockRoom,nil).Maybe()

		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil).Maybe()



		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		chatroom,err:=u.GetChatRoomsFor(context.TODO(),mockStudent.ID)
		assert.Error(t, err)
		assert.Nil(t, chatroom)

		mockRoomRepo.AssertExpectations(t)
	})
	t.Run("case error student does not exist or error in student repo", func(t *testing.T) {
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(nil,errors.New("error")).Once()



		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil).Maybe()
		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		chatroom,err:=u.GetChatRoomsFor(context.TODO(),mockStudent.ID)
		assert.Error(t, err)
		assert.Nil(t, chatroom)

		mockRoomRepo.AssertExpectations(t)
	})
	t.Run("case success", func(t *testing.T) {
		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil).Once()

		mockRoomRepo.On("GetRoomsFor",mock.Anything, mock.Anything).
			Return(&mockStudentChatRoom, nil).Once()

		mockRoomRepo.On("GetRoom",mock.Anything, mock.Anything).
			Return(&mockRoom,nil)

		mockStudentRepo.On("GetStudent",mock.Anything,mock.AnythingOfType("string")).
			Return(&mockStudent,nil)

		u := NewRoomUseCase(mockRoomRepo,mockStudentRepo,time.Second)
		chatroom,err:=u.GetChatRoomsFor(context.TODO(),mockStudent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, chatroom)

		mockRoomRepo.AssertExpectations(t)
	})

}

func TestDeleteRoom(t *testing.T) {

	mockRoomRepo := new(mocks.RoomRepository)
	mockStudentRepo := new(mocks.StudentRepository)
	var mockStudent domain.Student
	faker.FakeData(&mockStudent)
	var mockRoom domain.ChatRoom
	faker.FakeData(&mockRoom)
	t.Run("case error in the repo", func(t *testing.T) {
		mockRoom.Admin.ID=mockStudent.ID
		mockRoomRepo.On("GetRoom", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockRoom,nil)
		mockRoomRepo.On("RemoveRoomForParticipantsAndDeleteRoom", mock.Anything, &mockRoom).
			Return(errors.New("error"))

		u := NewRoomUseCase(mockRoomRepo, mockStudentRepo, time.Second)
		err := u.DeleteRoom(context.TODO(),  mockStudent.ID,mockRoom.RoomID)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})


	t.Run("case error the user trying to delete is not the admin", func(t *testing.T) {
		mockRoom.Admin.ID=mockStudent.ID+"error"
		mockRoomRepo.On("GetRoom", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockRoom,nil)


		u := NewRoomUseCase(mockRoomRepo, mockStudentRepo, time.Second)
		err := u.DeleteRoom(context.TODO(),  mockStudent.ID,mockRoom.RoomID)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})
	t.Run("case room does not exist", func(t *testing.T) {
		mockRoom.Admin.ID=mockStudent.ID
		mockRoomRepo.On("GetRoom", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil,errors.New("error"))


		u := NewRoomUseCase(mockRoomRepo, mockStudentRepo, time.Second)
		err := u.DeleteRoom(context.TODO(),  mockStudent.ID,mockRoom.RoomID)
		assert.Error(t, err)
		mockRoomRepo.AssertExpectations(t)

	})

	t.Run("case success", func(t *testing.T) {

		mockRoomRepo.On("GetRoom", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(&mockRoom,nil)
		mockRoomRepo.On("RemoveRoomForParticipantsAndDeleteRoom", mock.Anything, &mockRoom).
			Return(nil)
		mockRoom.Admin.ID=mockStudent.ID
		u := NewRoomUseCase(mockRoomRepo, mockStudentRepo, time.Second)
		err := u.DeleteRoom(context.TODO(),  mockStudent.ID,mockRoom.RoomID)
		assert.NoError(t, err)
		mockRoomRepo.AssertExpectations(t)

	})

}