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

// GetAllTags returns all the tags currently available
func (h *RoomHandler) SaveRoom(c *gin.Context) {



	var room domain.ChatRoom
	err := c.ShouldBindJSON(&room)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid data"))
		return
	}


	ctx := c.Request.Context()
	err = h.u.SaveRoom(ctx ,&room )
	if err != nil {
		switch v := err.(type) {
		case *errors.RestError:
			c.JSON(v.Code, v)
			return
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError(err.Error()))
			return
		}
	}




	c.JSON(http.StatusCreated, httputils.NewResponse("Room created"))
}