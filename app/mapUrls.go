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
		mh.ServeWs(c.Writer, c.Request, roomID, userID)
	})

	router.POST("/chat/room", rh.SaveRoom)
	router.GET("/chat/rooms/:id", rh.GetChatRoomsFor)
	router.PUT("/chat/add/:roomID/:id", rh.AddUserToRoom)
	router.PUT("/chat/remove/:roomID/:id", rh.RemoveUserFromRoom)
	router.DELETE("/chat/:id/:roomID", rh.DeleteRoom)
}
