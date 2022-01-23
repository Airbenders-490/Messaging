package app

import (
	"chat/messaging/delivery/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

func mapUrls(router *gin.Engine, mh *http.MessageHandler, rh *http.RoomHandler) {

	router.GET("/chat/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		// todo: get this from jwt token
		userID := strconv.Itoa(rand.Int())
		ctx := c.Request.Context()
		mh.ServeWs(c.Writer, c.Request, roomID, userID, ctx)
	})

	router.POST("/rooms", rh.SaveRoom)
	router.GET("/rooms/:id", rh.GetChatRoomsFor)
	router.PUT("/rooms/add/:roomID/:id", rh.AddUserToRoom)
	router.PUT("/rooms/remove/:roomID/:id", rh.RemoveUserFromRoom)
	router.DELETE("/rooms/:id/:roomID", rh.DeleteRoom)
}
