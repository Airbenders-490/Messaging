package http

import (
	"chat/domain"
	"chat/domain/mocks"
	"chat/utils/errors"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var mockRoom domain.ChatRoom
var mockRoomUseCase = new(mocks.RoomUseCase)
var rh = NewRoomHandler(mockRoomUseCase)
var contentType = "application/json"
const postRoomPath = "%s/rooms"
const readFailureMessage = "failed to read from message"

func TestSaveRoom(t *testing.T) {
	err := faker.FakeData(&mockRoom)
	assert.NoError(t, err)
	router := gin.Default()
	router.POST("/rooms", rh.SaveRoom)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("SaveRoom Success", func(t *testing.T) {
		mockRoomUseCase.
			On("SaveRoom", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		postBody, err := json.Marshal(mockRoom)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		response, err := server.Client().Post(fmt.Sprintf(postRoomPath, server.URL), contentType, reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})

	t.Run("Fail: Invalid request body format", func(t *testing.T) {
		response, err := server.Client().Post(fmt.Sprintf(postRoomPath, server.URL), contentType, nil)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
	})

	t.Run("Fail: SaveRoom error", func(t *testing.T) {
		mockRoomUseCase.
			On("SaveRoom", mock.Anything, mock.Anything).
			Return(errors.NewInternalServerError("")).
			Once()

		postBody, err := json.Marshal(mockRoom)
		assert.NoError(t, err)
		reader := strings.NewReader(string(postBody))
		response, err := server.Client().Post(fmt.Sprintf(postRoomPath, server.URL), contentType, reader)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})
}

func TestAddUserToRoom(t *testing.T) {
	router := gin.Default()
	router.PUT("/rooms/add/:roomID/:id", rh.AddUserToRoom)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("AddUserToRoom Success", func(t *testing.T) {
		mockRoomUseCase.
			On("AddUserToRoom", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/add/1/1", server.URL))
		request := http.Request{Method: "PUT", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusAccepted, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})

	t.Run("Fail: AddUserToRoom error", func(t *testing.T) {
		mockRoomUseCase.
			On("AddUserToRoom", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.NewInternalServerError("")).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/add/1/1", server.URL))
		request := http.Request{Method: "PUT", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})
}

func TestRemoveUserFromRoom(t *testing.T) {
	router := gin.Default()
	router.PUT("/rooms/remove/:roomID/:id", rh.RemoveUserFromRoom)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("RemoveUserFromRoom Success", func(t *testing.T) {
		mockRoomUseCase.
			On("RemoveUserFromRoom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/remove/1/1", server.URL))
		request := http.Request{Method: "PUT", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusAccepted, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})

	t.Run("Fail: RemoveUserFromRoom error", func(t *testing.T) {
		mockRoomUseCase.
			On("RemoveUserFromRoom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.NewInternalServerError("")).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/remove/1/1", server.URL))
		request := http.Request{Method: "PUT", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})
}

func TestGetChatRoomsFor(t *testing.T) {
	router := gin.Default()
	router.GET("/rooms/:id", rh.GetChatRoomsFor)
	server := httptest.NewServer(router)
	defer server.Close()
	studentChatRooms := &domain.StudentChatRooms{}

	t.Run("GetChatRoomsFor Success", func(t *testing.T) {
		mockRoomUseCase.
			On("GetChatRoomsFor", mock.Anything, mock.Anything).
			Return(studentChatRooms, nil).
			Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/rooms/1", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})

	t.Run("Fail: GetChatRoomsFor error", func(t *testing.T) {
		mockRoomUseCase.
			On("GetChatRoomsFor", mock.Anything, mock.Anything).
			Return(nil, errors.NewInternalServerError("")).
			Once()

		response, err := server.Client().Get(fmt.Sprintf("%s/rooms/1", server.URL))
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})
}

func TestDeleteRoom(t *testing.T) {
	router := gin.Default()
	router.DELETE("/rooms/:id/:roomID", rh.DeleteRoom)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("DeleteRoom Success", func(t *testing.T) {
		mockRoomUseCase.
			On("DeleteRoom", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/1/1", server.URL))
		request := http.Request{Method: "DELETE", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusAccepted, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})

	t.Run("Fail: DeleteRoom error", func(t *testing.T) {
		mockRoomUseCase.
			On("DeleteRoom", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.NewInternalServerError("")).
			Once()

		myUrl, err := url.Parse(fmt.Sprintf("%s/rooms/1/1", server.URL))
		request := http.Request{Method: "DELETE", URL: myUrl}
		response, err := server.Client().Do(&request)
		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			assert.Fail(t, readFailureMessage)
		}
		mockRoomUseCase.AssertExpectations(t)
	})
}
