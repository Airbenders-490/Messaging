package app

import (
	"chat/messaging/delivery/http"
	"chat/messaging/repository"
	"chat/messaging/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"log"
	"os"
	"time"
)

func Server(messageHandler *http.MessageHandler) *gin.Engine {
	router := gin.Default()
	mapUrls(router, messageHandler)
	return router
}

// Start runs the server
func Start() {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = os.Getenv("CASSANDRA_CHAT_KEYSPACE")

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected Cassandra database OK")

	mr := repository.NewChatRepository(session)
	rr := repository.NewRoomRepository(session)
	u := usecase.NewMessageUseCase(time.Second*2, mr, rr)
	messageHandler := http.NewMessageHandler(u)
	go http.MainHub.StartHubListener()
	router := Server(messageHandler)
	router.Run()
}
