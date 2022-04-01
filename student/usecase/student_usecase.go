package usecase

import (
	"chat/domain"
	"chat/student/repository"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const (
	failQueue = "Failed to declare a queue"
	failedToRegisterConsumer = "Failed to register a consumer"
	failedBindQueue = "Failed to bind a queue"
)

type studentUseCase struct {
	sr domain.StudentRepository
}

func NewStudentUseCase(repository repository.StudentRepository) domain.StudentUseCase {
	return &studentUseCase{sr: repository}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (s studentUseCase) ListenStudentCreation(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, failQueue)

	err = ch.QueueBind(
		q.Name,            // queue name
		"profile.created", // routing key
		"profile",         // exchange
		false,
		nil)
	failOnError(err, failedBindQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, failedToRegisterConsumer)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var st domain.Student
			json.Unmarshal(d.Body, &st)
			err = s.sr.SaveStudent(context.Background(), &st)
			if err != nil {
				log.Println("couldn't save student ", err)
				continue
			}
			d.Ack(false)
			log.Println("saved a student")
		}
	}()

	log.Printf(" [*] Waiting for student creation.")
	<-forever
}


// ListenStudentEdit listens to the
func (s studentUseCase) ListenStudentEdit(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, failQueue)

	err = ch.QueueBind(
		q.Name,            // queue name
		"profile.edited", // routing key
		"profile",         // exchange
		false,
		nil)
	failOnError(err, failedBindQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, failedToRegisterConsumer)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var st domain.Student
			json.Unmarshal(d.Body, &st)
			err = s.sr.EditStudent(context.Background(), &st)
			if err != nil {
				log.Println("couldn't edit student ", err)
				continue
			}
			d.Ack(false)
			log.Println("edited a student")
		}
	}()

	log.Printf(" [*] Waiting for student update.")
	<-forever
}

// ListenStudentDelete listens to the queue for any deleted students
func (s studentUseCase) ListenStudentDelete(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, failQueue)

	err = ch.QueueBind(
		q.Name,            // queue name
		"profile.deleted", // routing key
		"profile",         // exchange
		false,
		nil)
	failOnError(err, failedBindQueue)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, failedToRegisterConsumer)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			id := string(d.Body)
			err = s.sr.DeleteStudent(context.Background(), id)
			if err != nil {
				log.Println("couldn't delete student ", err)
				continue
			}
			d.Ack(false)
			log.Println("deleted a student")
		}
	}()

	log.Printf(" [*] Waiting for delete")
	<-forever
}