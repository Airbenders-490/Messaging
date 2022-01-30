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

func Server(mh *http.MessageHandler, rh *http.RoomHandler, mw Middleware) *gin.Engine {
	router := gin.Default()
	mapUrls(mw, router, mh, rh)
	return router
}

// Start runs the server
func Start() {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected Cassandra database OK")

	sr := repository.NewStudentRepository(session)
	su := usecase.NewStudentUseCase(*sr)

	go su.ListenStudentCreation()
	go su.ListenStudentEdit()
	go su.ListenStudentDelete()

	mr := repository.NewChatRepository(session)
	rr := repository.NewRoomRepository(session)

	mu := usecase.NewMessageUseCase(time.Second*2, mr, rr)
	mh := http.NewMessageHandler(mu)

	ru := usecase.NewRoomUseCase(rr, sr, time.Second*2)
	rh := http.NewRoomHandler(ru)

	mw := NewMiddleware()

	go http.MainHub.StartHubListener()
	router := Server(mh, rh, mw)
	router.Run()
}
