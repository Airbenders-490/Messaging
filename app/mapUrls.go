package app

import (
	"chat/messaging/delivery/http"
	roomHttp "chat/room/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapChatUrls(mw Middleware, r *gin.Engine, mh *http.MessageHandler) {

	r.GET("/chat/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		ctx := c.Request.Context()
		mh.ServeWs(c.Writer, c.Request, roomID, ctx)
	})
	router := r.Group("/api")
	router.Use(mw.AuthMiddleware())

	const pathRoomID = "chat/:roomID"
	router.POST(pathRoomID, mh.LoadMessages)
	router.PUT(pathRoomID, mh.EditMessage)
	router.DELETE(pathRoomID, mh.DeleteMessage)
	router.POST("chat/joinRequest/:roomID", mh.JoinRequest)
	router.POST("chat/rejectRequest/:roomID/:userID", mh.RejectJoinRequest)
}

func mapRoomURLs(mw Middleware, r *gin.Engine, rh *roomHttp.RoomHandler) {
	router := r.Group("/api/rooms")
	router.Use(mw.AuthMiddleware())

	router.POST("", rh.SaveRoom)
	router.GET("", rh.GetChatRoomsFor)
	router.GET("/class/:className", rh.GetChatRoomsByClass)
	router.PUT("/add/:roomID/:id", rh.AddUserToRoom)
	router.PUT("/remove/:roomID/:id", rh.RemoveUserFromRoom)
	router.DELETE("/:roomID", rh.DeleteRoom)
}