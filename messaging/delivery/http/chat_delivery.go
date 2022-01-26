package http

import (
	"chat/domain"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws   *websocket.Conn
	send chan domain.Message
}

type subscription struct {
	conn   *connection
	roomID string
	userID string
}

// Hub is the heart of the chat app. This is what is used to hold "rooms", register and unregister when connecting and
// disconnecting, and broadcast. Whenever a message is sent to broadcast channel, it is delivered to all the connections
// in room
type Hub struct {
	rooms      map[string]map[subscription]bool
	broadcast  chan domain.Message
	Register   chan subscription
	unregister chan subscription
}

// MainHub is the instantiation of our chat's heart
// todo: make this a singleton
var MainHub = Hub{
	broadcast:  make(chan domain.Message),
	Register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[subscription]bool),
}

const (
	maxMessageSize int64 = 1024
	pongWait             = time.Minute * 5
	pingPeriod           = time.Minute * 5
	writeWait            = time.Minute
)

func (s *subscription) readPump(u domain.MessageUseCase) {
	c := s.conn
	defer func() {
		MainHub.unregister <- *s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	_ = c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { _ = c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := domain.Message{RoomID: s.roomID, SentTimestamp: time.Now(), FromStudentID: s.userID, MessageBody: string(msg)}
		err = u.SaveMessage(context.Background(), &m)
		if err != nil {
			log.Printf("Failed to save message with err %s", err.Error())
		}
		MainHub.broadcast <- m
	}
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			res, err := json.Marshal(message)
			if err != nil {
				log.Printf("message %s couldn't be sent to %s in room %s.", message, s.userID, s.roomID)
			}
			if err = c.write(websocket.TextMessage, res); err != nil {
				log.Printf("message %s couldn't be sent to %s in room %s.", message, s.userID, s.roomID)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// ServeWs is the handleFunc for connecting to a room's websocket. A user must be authorized, i.e. already added to
// the room before he can connect to the room. Otherwise, returns 401.
func (h *MessageHandler) ServeWs(w http.ResponseWriter, r *http.Request, roomID string, userID string, ctx context.Context) {
	authorized := h.u.IsAuthorized(ctx, userID, roomID)
	if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "Not authorized to enter the room number "+roomID)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan domain.Message), ws: ws}
	s := subscription{c, roomID, userID}
	if authorized {
		MainHub.Register <- s
	}
	go s.writePump()
	go s.readPump(h.u)
}

// StartHubListener starts the functioning of the Hub with respect to registering, unregistering and broadcasting
func (h *Hub) StartHubListener() {
	for {
		select {
		case s := <-h.Register:
			connections := h.rooms[s.roomID]
			if connections == nil {
				connections = make(map[subscription]bool)
				h.rooms[s.roomID] = connections
			}
			h.rooms[s.roomID][s] = true
		case s := <-h.unregister:
			connections := h.rooms[s.roomID]
			if connections != nil {
				if _, ok := connections[s]; ok {
					delete(connections, s)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.roomID)
					}
				}
			}
		case m := <-h.broadcast:
			subscriptions := h.rooms[m.RoomID]
			for s := range subscriptions {
				if m.FromStudentID == s.userID {
					continue
				}
				select {
				case s.conn.send <- m:
				default:
					close(s.conn.send)
					delete(subscriptions, s)
					if len(subscriptions) == 0 {
						delete(h.rooms, m.RoomID)
					}
				}
			}

		}
	}
}

// MessageHandler is the standard delivery handler for messaging service
type MessageHandler struct {
	u domain.MessageUseCase
}

// NewMessageHandler instantiates and returns a new MessageHandler
func NewMessageHandler(u domain.MessageUseCase) *MessageHandler {
	return &MessageHandler{u: u}
}
