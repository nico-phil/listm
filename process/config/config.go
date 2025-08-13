package config

import (
	"os"
	"strings"
)

func GetContactPoints() []string {
	valueStr := os.Getenv("CASSANDRA_CONTACT_POINTS")
	if valueStr == "" {
		return []string{"cassandra-node1:9042"}
	}

	valueArray := strings.Split(valueStr, "")
	return valueArray
}

func GetKeyspace() string {
	keypsace := os.Getenv("CASSANDRA_KEYSPACE")
	if keypsace == "" {
		return "pluralistmanagement"
	}

	return keypsace
}
