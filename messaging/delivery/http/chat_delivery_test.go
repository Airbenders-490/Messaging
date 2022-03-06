package http_test

import (
	"chat/app"
	"chat/domain"
	"chat/domain/mocks"
	"chat/messaging/delivery/http"
	"chat/utils/errors"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const chatRoomAddr = "%s/chat/%s"
const messageBody = "Hello sir!"
const putChatPath = "/api/chat/%s"
const invalidDataMessage = "invalid data"
const invalidBodyMessage = "invalid body"
const restError = "rest error"
const errorOccurredMessage = "error occurred"
const validChatRoomID = "1"
const invalidChatRoomID = "2"
const errorMassage = "can't write message to socket"

func TestMessageSending(t *testing.T) {
	mockMessageUsecase := new(mocks.MessageUseCase)
	mh := http.NewMessageHandler(mockMessageUsecase)
	mw := new(mocks.MiddlewareMock)
	server := httptest.NewServer(app.Server(mh, nil, mw))
	defer server.Close()
	mockMessageUsecase.
		On("SaveMessage", mock.Anything, mock.Anything).Return(nil)

	addr, err := url.Parse(server.URL)
	if err != nil {
		assert.Fail(t, "unable to get test server url")
	}
	addr.Scheme = "ws"
	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	const validChatRoomID = "1"
	const invalidChatRoomID = "2"
	t.Run("success", func(t *testing.T) {
		// return only twice since connecting twice
		mockMessageUsecase.
			On("IsAuthorized", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true).Twice()
		// establish a connection. This is our "monitor" connection. We will send messages to this chat-room and
		// monitor this to read a valid response
		wsDefault, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		response, errChan := readyToReadMethod(wsDefault)

		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, errorMassage)
		select {
		case r := <-response:
			var event http.Event
			err = json.Unmarshal(r, &event)
			assert.NoError(t, err, "error unmarshalling")
			assert.Equal(t, messageBody, event.Message.MessageBody, "invalid message body")
			assert.Equal(t, validChatRoomID, event.Message.RoomID, "invalid room id received")
			// todo: test the id after auth connection
		case e := <-errChan:
			assert.Fail(t, e.Error())
		}
	})

	// test if unauthorized person can register
	t.Run("unauthorized: fail to handshake", func(t *testing.T) {
		// mock authorized once to register, and false second time to not register
		mockMessageUsecase.
			On("IsAuthorized", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(false).Once()
		// establish a connection. This is our "monitor" connection. We will send messages to this chat-room and
		// monitor this to read a valid response
		_, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		assert.Error(t, err)
	})

	// unregister
	t.Run("unregister", func(t *testing.T) {
		// return only twice since connecting twice
		mockMessageUsecase.
			On("IsAuthorized", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true).Twice()
		// establish a connection. This is our "monitor" connection. We will send messages to this chat-room and
		// monitor this to read a valid response. In this case, nothing should be received after closing the conn
		// register the reader and then unregister
		wsDefault, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		_ = wsDefault.Close()

		ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		response, errChan := readyToReadMethod(wsDefault)

		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, errorMassage)
		select {
		case <-response:
			assert.Fail(t, "received a message when not expected")
		case e := <-errChan:
			assert.Error(t, e)
		}
	})

	t.Run("different room don't receive", func(t *testing.T) {
		// return only twice since connecting twice
		mockMessageUsecase.
			On("IsAuthorized", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true).Twice()
		// establish a connection. This is our "monitor" connection. We will send messages to this chat-room and
		// monitor this to read a valid response
		wsDefault, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), validChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(chatRoomAddr, addr.String(), invalidChatRoomID), nil)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		response, errChan := readyToReadMethod(wsDefault)

		_ = wsDefault.SetReadDeadline(time.Now().Add(time.Second))
		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, errorMassage)
		select {
		case <-response:
			assert.Fail(t, "received message in a different room. What?!")
		case e := <-errChan:
			assert.Error(t, e)
		}
	})
}

func readyToReadMethod(wsDefault *websocket.Conn) (chan []byte, chan error) {
	response := make(chan []byte)
	errChan := make(chan error)
	readyToRead := make(chan bool)
	go func(r chan []byte, e chan error) {
		readyToRead <- true
		_ = wsDefault.SetReadDeadline(time.Now().Add(time.Second))
		_, msg, er := wsDefault.ReadMessage()
		if er != nil {
			e <- er
		} else {
			r <- msg
		}
	}(response, errChan)
	<-readyToRead
	return response, errChan
}

func TestLoadMessages(t *testing.T) {
	mockUseCase := new(mocks.MessageUseCase)
	mh := http.NewMessageHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	r := app.Server(mh, nil, mw)

	var retrievedMessages []domain.Message
	err := faker.FakeData(&retrievedMessages)
	assert.NoError(t, err)
	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	var mockMessage domain.Message
	err = faker.FakeData(&mockMessage.SentTimestamp)
	assert.NoError(t, err)
	err = faker.FakeData(&mockMessage.RoomID)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockUseCase.On("GetMessages", mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("int")).
			Return(retrievedMessages, nil).
			Once()
		getBody, err := json.Marshal(mockMessage)
		assert.NoError(t, err)
		reader := strings.NewReader(string(getBody))
		reqFound := httptest.NewRequest("POST", fmt.Sprintf("/api/chat/%s?limit=%s", mockMessage.RoomID, "5"), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("no limit", func(t *testing.T) {
		mockUseCase.On("GetMessages", mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("int")).
			Return(retrievedMessages, nil).
			Once()
		getBody, err := json.Marshal(mockMessage)
		assert.NoError(t, err)
		reader := strings.NewReader(string(getBody))
		reqFound := httptest.NewRequest("POST", fmt.Sprintf(putChatPath, mockMessage.RoomID), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run(invalidDataMessage, func(t *testing.T) {
		reader := strings.NewReader(invalidBodyMessage)
		reqFound := httptest.NewRequest("POST", fmt.Sprintf(putChatPath, mockMessage.RoomID), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})


	t.Run(restError, func(t *testing.T) {
		restErr := errors.NewConflictError(errorOccurredMessage)
		mockUseCase.On("GetMessages", mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("int")).
			Return(nil, restErr).
			Once()
		getBody, err := json.Marshal(mockMessage)
		assert.NoError(t, err)
		reader := strings.NewReader(string(getBody))
		reqFound := httptest.NewRequest("POST", fmt.Sprintf("/api/chat/%s?limit=%s", mockMessage.RoomID, "5"), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestEditMessage(t *testing.T) {
	mockUseCase := new(mocks.MessageUseCase)
	mh := http.NewMessageHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	r := app.Server(mh, nil, mw)

	var editedMessage domain.Message
	err := faker.FakeData(&editedMessage)
	assert.NoError(t, err)

	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	t.Run("success", func(t *testing.T) {
		putBody, err := json.Marshal(editedMessage)
		assert.NoError(t, err)
		mockUseCase.On("EditMessage", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("string")).
			Return(&editedMessage, nil).Once()

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putChatPath,
			editedMessage.RoomID), reader)

		reqFound.Header.Set("id", editedMessage.FromStudentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 200, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run(invalidDataMessage, func(t *testing.T) {
		reader := strings.NewReader(invalidBodyMessage)
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putChatPath, editedMessage.RoomID), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("unauthorized user", func(t *testing.T) {
		putBody, err := json.Marshal(editedMessage)
		assert.NoError(t, err)

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putChatPath,
			editedMessage.RoomID), reader)

		reqFound.Header.Set("id", "avc")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 401, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run(restError, func(t *testing.T) {
		putBody, err := json.Marshal(editedMessage)
		assert.NoError(t, err)
		restErr := errors.NewConflictError(errorOccurredMessage)
		mockUseCase.On("EditMessage", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("string")).
			Return(nil, restErr).Once()

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("PUT", fmt.Sprintf(putChatPath,
			editedMessage.RoomID), reader)

		reqFound.Header.Set("id", editedMessage.FromStudentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDeleteMessage(t *testing.T) {
	mockUseCase := new(mocks.MessageUseCase)
	mh := http.NewMessageHandler(mockUseCase)
	mw := new(mocks.MiddlewareMock)
	r := app.Server(mh, nil, mw)

	var deletedMessage domain.Message
	err := faker.FakeData(&deletedMessage)
	assert.NoError(t, err)
	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	t.Run("success", func(t *testing.T) {
		putBody, err := json.Marshal(deletedMessage)
		assert.NoError(t, err)
		mockUseCase.On("DeleteMessage", mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string")).
			Return(nil).Once()

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putChatPath,
			deletedMessage.RoomID), reader)

		reqFound.Header.Set("id", deletedMessage.FromStudentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 202, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run(invalidDataMessage, func(t *testing.T) {
		//var mockMessage domain.Message
		reader := strings.NewReader(invalidBodyMessage)
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putChatPath, deletedMessage.RoomID), reader)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 400, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("unauthorized user", func(t *testing.T) {
		putBody, err := json.Marshal(deletedMessage)
		assert.NoError(t, err)

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putChatPath,
			deletedMessage.RoomID), reader)

		reqFound.Header.Set("id", "avc")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, 401, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run(restError, func(t *testing.T) {
		putBody, err := json.Marshal(deletedMessage)
		assert.NoError(t, err)
		restErr := errors.NewConflictError(errorOccurredMessage)
		mockUseCase.On("DeleteMessage", mock.Anything, mock.AnythingOfType("string"),
			mock.Anything, mock.AnythingOfType("string")).Return(restErr).Once()

		reader := strings.NewReader(string(putBody))
		reqFound := httptest.NewRequest("DELETE", fmt.Sprintf(putChatPath,
			deletedMessage.RoomID), reader)

		reqFound.Header.Set("id", deletedMessage.FromStudentID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, reqFound)
		assert.Equal(t, restErr.Code, w.Code)
		mockUseCase.AssertExpectations(t)
	})

}
