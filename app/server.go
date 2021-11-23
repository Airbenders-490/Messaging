package app

import (
	"chat/repository"
	"github.com/gocql/gocql"
	"log"
	"os"
)


// Start runs the server
func Start() {

	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = os.Getenv("CASSANDRA_CHAT_KEYSPACE")

	session, err := cluster.CreateSession()
	if  err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected Cassandra database OK")

	chatRepository := repository.NewChatRepository(session)
	chatRepository.GetByFromAndToID("40025120", "400000") // example, can remove

	// TODO: add middleware for authentication
	// TODO: map URLs & add router
}
