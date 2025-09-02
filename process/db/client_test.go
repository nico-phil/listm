package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSession tests db session
func TestGetSession(t *testing.T) {
	session = nil
	result := Getsession()
	assert.Nil(t, result)
}

// TestCloseSession tests close db session
func TestCloseSession(t *testing.T) {
	// Test closing when session is nil
	session = nil

	//should not panic
	CloseSession()

	// Test is primarily to ensure no panic accurs
	assert.True(t, true)
}
