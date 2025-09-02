package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContactPoints(t *testing.T) {

	var tests = []struct {
		name     string
		envValue string
		expected []string
	}{
		{
			name:     "default contact points",
			envValue: "",
			expected: []string{"cassandra-node1:9042"},
		},

		{
			name:     "single contact points from env",
			envValue: "cassandra-node1:9042",
			expected: []string{"cassandra-node1:9042"},
		},

		{
			name:     "multiple contact points",
			envValue: "node1:9042 node2:9042 node3:9042",
			expected: []string{"node1:9042", "node2:9042", "node3:9042"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CASSANDRA_CONTACT_POINTS", tt.envValue)
			result := GetContactPoints()
			assert.Equal(t, tt.expected, result)
		})
	}

	// clean up
	os.Unsetenv("CASSANDRA_CONTACT_POINTS")
}

func TestGetKeyspace(t *testing.T) {
	cases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "default keyspace",
			envValue: "",
			expected: "pluralistmanagement",
		},

		{
			name:     "keyspace from env",
			envValue: "mykeyspace",
			expected: "mykeyspace",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Setenv("CASSANDRA_KEYSPACE", c.envValue)
			result := GetKeyspace()
			assert.Equal(t, c.expected, result)
		})
	}

	os.Unsetenv("CASSANDRA_KEYSPACE")
}

func TestGetRedisAddr(t *testing.T) {
	cases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "default redis addr",
			envValue: "",
			expected: "localhost:6379",
		},

		{
			name:     "redis addr from env",
			envValue: "my_redis_addr:6379",
			expected: "my_redis_addr:6379",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Setenv("REDIS_ADDR", c.envValue)
			result := GetRedisArr()
			assert.Equal(t, c.expected, result)
		})
	}

	// clear env
	os.Unsetenv("REDIS_ADDR")
}

func TestGetRedisPASSWORD(t *testing.T) {
	cases := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "default redis password",
			envValue: "",
			expected: "",
		},

		{
			name:     "redis password from env",
			envValue: "redis_password",
			expected: "redis_password",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Setenv("REDIS_PASSWORD", c.envValue)
			result := GetRedisPassword()
			assert.Equal(t, c.expected, result)
		})
	}

	// clear env
	os.Unsetenv("REDIS_PASSWORD")
}
