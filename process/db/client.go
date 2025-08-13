package db

import (
	"errors"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/nico-phil/process/config"
)

var (
	ErrNoConnection = errors.New("no database connection ")
)

var session *gocql.Session

func NewClient() error {

	contactPoints := config.GetContactPoints()
	keypsace := config.GetKeyspace()

	cluster := gocql.NewCluster(contactPoints...)
	cluster.Keyspace = keypsace
	cluster.Port = 9142
	cluster.Timeout = 10 * time.Second

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Printf("failed to connect to cassandra %v", err)
		return err
	}
	return nil
}

func Getsession() *gocql.Session {
	return session
}
