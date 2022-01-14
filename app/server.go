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

func Server(mh *http.MessageHandler, rh *http.RoomHandler) *gin.Engine {
	router := gin.Default()
	mapUrls(router, mh, rh)
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
	sr := repository.NewStudentRepository(session)

	mu := usecase.NewMessageUseCase(time.Second*2, mr, rr)
	mh := http.NewMessageHandler(mu)

	ru := usecase.NewRoomUseCase(rr, sr, time.Second*2)
	rh := http.NewRoomHandler(ru)

	go http.MainHub.StartHubListener()
	router := Server(mh, rh)
	router.Run()
}
