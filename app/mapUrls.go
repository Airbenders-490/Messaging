package app

import (
	"chat/messaging/delivery/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

<<<<<<< refs/remotes/origin/resolve-techinical-debt/STUD-234
const pathRoomID = "api/chat/:roomID"
func mapUrls(mw Middleware, router *gin.Engine, mh *http.MessageHandler, rh *http.RoomHandler) {
=======
func mapUrls(mw Middleware, r *gin.Engine, mh *http.MessageHandler, rh *http.RoomHandler) {
>>>>>>> STUD248/app: fix auth for chatroom connection

	r.GET("/chat/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		userID := strconv.Itoa(rand.Int())
		ctx := c.Request.Context()
		mh.ServeWs(c.Writer, c.Request, roomID, userID, ctx)
	})

	router := r.Group("/api")
	router.Use(mw.AuthMiddleware())

	router.POST("/rooms", rh.SaveRoom)
	router.GET("/rooms/:id", rh.GetChatRoomsFor)
	router.PUT("/rooms/add/:roomID/:id", rh.AddUserToRoom)
	router.PUT("/rooms/remove/:roomID/:id", rh.RemoveUserFromRoom)
	router.DELETE("/rooms/:id/:roomID", rh.DeleteRoom)

<<<<<<< refs/remotes/origin/resolve-techinical-debt/STUD-234

	router.GET(pathRoomID, mh.LoadMessages)
	router.PUT(pathRoomID, mh.EditMessage)
	router.DELETE(pathRoomID, mh.DeleteMessage)
=======
	router.GET("chat/:roomID", mh.LoadMessages)
	router.PUT("chat/:roomID", mh.EditMessage)
	router.DELETE("chat/:roomID", mh.DeleteMessage)
>>>>>>> STUD248/app: fix auth for chatroom connection
}
