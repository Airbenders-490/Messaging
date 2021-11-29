package app

import (
	"chat/messaging/delivery/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

func mapUrls(router *gin.Engine, mh *http.MessageHandler) {
	router.GET("/chat/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		// todo: get this from jwt token
		userID := strconv.Itoa(rand.Int())
		mh.ServeWs(c.Writer, c.Request, roomID, userID)
	})
}
