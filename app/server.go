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
	"github.com/streadway/amqp"
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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

	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"profile", // name
		"topic",   // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	go su.ListenStudentCreation(ch)
	go su.ListenStudentEdit(ch)
	go su.ListenStudentDelete(ch)

	mw := NewMiddleware()

	mainHub := http.NewHub()
	go mainHub.StartHubListener()
	router := Server(mh, rh, mw)
	router.Run()
}
