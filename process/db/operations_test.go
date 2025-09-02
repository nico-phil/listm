package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetCampaigns_NoConnection test getcampaign when no session
func TestGetCampaigns_NoConnection(t *testing.T) {

	// save original session
	originalsession := session
	defer func() {
		session = originalsession
	}()

	// set the session to nil, since we are connected to db
	session = nil
	_, err := GetAllCampaigns()
	assert.Nil(t, session)
	assert.Equal(t, ErrNoConnection, err)
}

func TestGetLists_NoConnection(t *testing.T) {
	originalsession := session
	defer func() {
		session = originalsession
	}()

	session = nil
	_, err := GetAllLists()
	assert.Nil(t, session)
	assert.Equal(t, ErrNoConnection, err)
}
