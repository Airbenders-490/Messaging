package app

import (
	"chat/messaging/delivery/http"
	"chat/messaging/repository"
	"chat/messaging/repository/cassandra"
	"chat/messaging/usecase"
	"chat/utils"
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

	mail := utils.NewSimpleMail()
	mr := repository.NewChatRepository(cassandra.NewSession(session))
	rr := repository.NewRoomRepository(cassandra.NewSession(session))
	sr := repository.NewStudentRepository(cassandra.NewSession(session))

	mu := usecase.NewMessageUseCase(time.Second*2, mr, rr, sr, mail)
	ru := usecase.NewRoomUseCase(rr, sr, time.Second*2)

	mh := http.NewMessageHandler(mu)
	rh := http.NewRoomHandler(ru)

	su := usecase.NewStudentUseCase(*sr)

	go su.ListenStudentCreation()
	go su.ListenStudentEdit()
	go su.ListenStudentDelete()

	mw := NewMiddleware()

	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	router := Server(mh, rh, mw)
	router.Run()
}
