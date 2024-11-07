package db

import (
	"github.com/gocql/gocql"
	"log"
)

var Session *gocql.Session

func InitCassandra() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1:9042", "127.0.0.1:9043")
	cluster.Keyspace = "recommendation_keyspace"
	cluster.Consistency = gocql.Quorum

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Erro ao conectar ao Cassandra:", err)
	}
	return Session
}

func CloseCassandra() {
	Session.Close()
}
