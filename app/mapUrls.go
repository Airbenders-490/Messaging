package app

import (
	"chat/messaging/delivery/http"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine, mh *http.MessageHandler) {
	router.GET("/chat/:roomID", mh.JoinChatRoom)
}
