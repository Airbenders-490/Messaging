package http

import (
	"chat/domain"
	"chat/utils/errors"
	"chat/utils/httputils"

	"github.com/gin-gonic/gin"
	"net/http"
)

type RoomHandler struct {
	u domain.RoomUseCase
}

// NewReviewHandler is the constructor
func NewRoomHandler(ru domain.RoomUseCase) *RoomHandler {
	return &RoomHandler{u: ru}
}

func (h *RoomHandler) SaveRoom(c *gin.Context) {
	var room domain.ChatRoom
	err := c.ShouldBindJSON(&room)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("Invalid request body format for Chat Room"))
		return
	}

	ctx := c.Request.Context()
	err = h.u.SaveRoom(ctx, &room)
	if err != nil {
		errors.SetRESTError(err, c)
		return
	}

	c.JSON(http.StatusCreated, httputils.NewResponse("Room created"))
}

func (h *RoomHandler) AddUserToRoom(c *gin.Context) {
	userID := c.Params.ByName("id")
	roomID := c.Params.ByName("roomID")
	if userID == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("userID AND roomID parameters must be provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.u.AddUserToRoom(ctx, roomID, userID)
	if err != nil {
		errors.SetRESTError(err, c)
		return
	}

	c.JSON(http.StatusAccepted, httputils.NewResponse("User added to Room"))
}

func (h *RoomHandler) RemoveUserFromRoom(c *gin.Context) {
	userID := c.Params.ByName("id")
	roomID := c.Params.ByName("roomID")
	if userID == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("userID AND roomID parameters must be provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.u.RemoveUserFromRoom(ctx, roomID, userID)
	if err != nil {
		errors.SetRESTError(err, c)
		return
	}

	c.JSON(http.StatusAccepted, httputils.NewResponse("User Removed from Room"))
}

func (h *RoomHandler) GetChatRoomsFor(c *gin.Context) {
	userID := c.Params.ByName("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("userID parameter must be provided"))
		return
	}

	ctx := c.Request.Context()
	studentChatRooms, err := h.u.GetChatRoomsFor(ctx, userID)
	if err != nil {
		errors.SetRESTError(err, c)
		return
	}

	c.JSON(http.StatusOK, studentChatRooms)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userID := c.Params.ByName("id")
	roomID := c.Params.ByName("roomID")
	if userID == "" || roomID == "" {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("userID AND roomID parameters must be provided"))
		return
	}

	ctx := c.Request.Context()
	err := h.u.DeleteRoom(ctx, userID, roomID)
	if err != nil {
		errors.SetRESTError(err, c)
		return
	}

	c.JSON(http.StatusAccepted, httputils.NewResponse("Room Deleted"))
}
