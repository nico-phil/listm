package tz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewZipcodeCache(t *testing.T) {
	zipCodeCache := NewZipcodeCache()

	assert.NotNil(t, zipCodeCache)
	assert.NotNil(t, zipCodeCache.cache)
}

// TestSetZipCode tests the Set zipcode function
func TestSetZipCode(t *testing.T) {
	zipCodeCache := NewZipcodeCache()

	zipCodeInfo := ZipCodeInfo{
		Zipcode:   "1234",
		Latitude:  40.7128,
		Longitude: -74.0060,
		City:      "New York",
		State:     "New York",
		TimeZone:  "America/New_York",
	}

	zipCodeCache.Set(zipCodeInfo.Zipcode, &zipCodeInfo)

	insertedZipCodeInfo, ok := zipCodeCache.cache[zipCodeInfo.Zipcode]

	// verify the zipcode in found
	assert.True(t, ok)
	assert.Equal(t, zipCodeInfo.Zipcode, insertedZipCodeInfo.Zipcode)
	assert.Equal(t, zipCodeInfo.Latitude, insertedZipCodeInfo.Latitude)
	assert.Equal(t, zipCodeInfo.Longitude, insertedZipCodeInfo.Longitude)
	assert.Equal(t, zipCodeInfo.City, insertedZipCodeInfo.City)
	assert.Equal(t, zipCodeInfo.State, insertedZipCodeInfo.State)
	assert.Equal(t, zipCodeInfo.TimeZone, insertedZipCodeInfo.TimeZone)

	_, ok = zipCodeCache.cache["4444"]
	assert.False(t, ok)

}
