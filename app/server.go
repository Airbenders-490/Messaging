package app

import (
	"chat/messaging/delivery/http"
	"chat/messaging/repository"
	"chat/messaging/repository/cassandra"
	"chat/messaging/usecase"
	http2 "chat/room/delivery/http"
	roomRepository "chat/room/repository"
	roomUseCase "chat/room/usecase"
	studentRepository "chat/student/repository"
	studentUseCase "chat/student/usecase"
	"chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"log"
	"os"
	"time"
)

func Server(mh *http.MessageHandler, rh *http2.RoomHandler, mw Middleware) *gin.Engine {
	router := gin.Default()
	mapChatUrls(mw, router, mh)
	mapRoomURLs(mw, router, rh)
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
	rr := roomRepository.NewRoomRepository(cassandra.NewSession(session))
	sr := studentRepository.NewStudentRepository(cassandra.NewSession(session))

	mu := usecase.NewMessageUseCase(time.Second*2, mr, rr, sr, mail)
	ru := roomUseCase.NewRoomUseCase(rr, sr, time.Second*2)

	mh := http.NewMessageHandler(mu)
	rh := http2.NewRoomHandler(ru)

	su := studentUseCase.NewStudentUseCase(*sr)

	go su.ListenStudentCreation()
	go su.ListenStudentEdit()
	go su.ListenStudentDelete()

	mw := NewMiddleware()

	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	router := Server(mh, rh, mw)
	router.Run()
}
