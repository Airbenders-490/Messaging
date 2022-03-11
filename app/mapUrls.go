package app

import (
	"chat/messaging/delivery/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

func mapUrls(mw Middleware, r *gin.Engine, mh *http.MessageHandler, rh *http.RoomHandler) {

	r.GET("/chat/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		userID := strconv.Itoa(rand.Int())
		ctx := c.Request.Context()
		mh.ServeWs(c.Writer, c.Request, roomID, userID, ctx)
	})
	router := r.Group("/api")
	router.Use(mw.AuthMiddleware())

	router.POST("/rooms", rh.SaveRoom)
	router.GET("/rooms", rh.GetChatRoomsFor)
	router.GET("/rooms/class/:className", rh.GetChatRoomsByClass)
	router.PUT("/rooms/add/:roomID/:id", rh.AddUserToRoom)
	router.PUT("/rooms/remove/:roomID/:id", rh.RemoveUserFromRoom)
	router.DELETE("/rooms/:roomID", rh.DeleteRoom)

	const pathRoomID = "chat/:roomID"
	router.POST(pathRoomID, mh.LoadMessages)
	router.PUT(pathRoomID, mh.EditMessage)
	router.DELETE(pathRoomID, mh.DeleteMessage)
	router.POST("chat/joinRequest/:roomID", mh.JoinRequest)
	router.POST("chat/rejectRequest/:roomID/:userID", mh.RejectRequest)
}
