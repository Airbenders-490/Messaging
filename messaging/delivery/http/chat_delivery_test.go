package http_test

import (
	"chat/app"
	"chat/domain/mocks"
	"chat/messaging/delivery/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const chatRoomAddr = "%s/chat/%s"

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
	go http.MainHub.StartHubListener()
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

		response := make(chan []byte)
		errChan := make(chan error)
		readyToRead := make(chan bool)
		const messageBody = "Hello sir!"
		go func(r chan []byte, e chan error) {
			readyToRead <- true
			_, msg, er := wsDefault.ReadMessage()
			if er != nil {
				e <- er
			} else {
				r <- msg
			}
		}(response, errChan)
		<- readyToRead
		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, "can't write message to socket")
		select {
		case r := <- response:
			var event http.Event
			err = json.Unmarshal(r, &event)
			assert.NoError(t, err, "error unmarshalling")
			assert.Equal(t, messageBody, event.Message.MessageBody, "invalid message body")
			assert.Equal(t, validChatRoomID, event.Message.RoomID, "invalid room id received")
			// todo: test the id after auth connection
		case e := <- errChan:
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

		response := make(chan []byte)
		errChan := make(chan error)
		readyToRead := make(chan bool)
		const messageBody = "Hello sir!"
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
		<- readyToRead
		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, "can't write message to socket")
		select {
		case <- response:
			assert.Fail(t, "received a message when not expected")
		case e := <- errChan:
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

		response := make(chan []byte)
		errChan := make(chan error)
		readyToRead := make(chan bool)
		const messageBody = "Hello sir!"
		go func(r chan []byte, e chan error) {
			readyToRead <- true
			_, msg, er := wsDefault.ReadMessage()
			if er != nil {
				e <- er
			} else {
				r <- msg
			}
		}(response, errChan)
		<- readyToRead
		_ = wsDefault.SetReadDeadline(time.Now().Add(time.Second))
		err = ws.WriteMessage(websocket.TextMessage, []byte(messageBody))
		assert.NoError(t, err, "can't write message to socket")
		select {
		case <- response:
			assert.Fail(t, "received message in a different room. What?!")
		case e := <- errChan:
			assert.Error(t, e)
		}
	})
}
