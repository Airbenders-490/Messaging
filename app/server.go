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

	repository.NewChatRepository(session)
}
