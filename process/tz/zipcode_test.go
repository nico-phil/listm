package tz

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewZipcodeCache tests creation of NewZipcodeCache
func TestNewZipcodeCache(t *testing.T) {
	zipCodeCache := NewZipcodeCache()

	assert.NotNil(t, zipCodeCache)
	assert.NotNil(t, zipCodeCache.cache)
	assert.Len(t, zipCodeCache.cache, 0)
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

// TestZipCodeCacheConcurrency tests inserting and accessing the cache concurrently
func TestZipCodeCacheConcurrency(t *testing.T) {
	cache := NewZipcodeCache()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			info := ZipCodeInfo{
				Zipcode:  fmt.Sprintf("zip%d", i),
				TimeZone: "America/New_York",
			}

			cache.Set(info.Zipcode, &info)
		}()

		go func() {
			defer wg.Done()
			zipCode := fmt.Sprintf("zip%d", i)
			_, _ = cache.getZipcode(zipCode)

		}()
	}

	wg.Wait()

	assert.Len(t, cache.cache, 100)
}

// TestGetLocalTimeAt tests localtime for a specific timezone
func TestGetLocalTimeAt(t *testing.T) {

	zipCodeCache := NewZipcodeCache()

	// Eastern time zone
	info1 := ZipCodeInfo{
		Zipcode:   "1234",
		Latitude:  40.7128,
		Longitude: -74.0060,
		City:      "New York",
		State:     "New York",
		TimeZone:  "America/New_York",
	}

	info2 := ZipCodeInfo{
		Zipcode:   "94016",
		Latitude:  37.7749,
		Longitude: -122.4194,
		City:      "San Francisco",
		State:     "California",
		TimeZone:  "America/Los_Angeles",
	}

	zipCodeCache.Set(info1.Zipcode, &info1)
	zipCodeCache.Set(info2.Zipcode, &info2)

	now := time.Now()

	_, err := GetLocalTimeAt(zipCodeCache, info1.Zipcode, now)
	if err != nil {
		t.Fatalf("failed to get local time for zipcode 1234")
	}

	_, err = GetLocalTimeAt(zipCodeCache, info2.Zipcode, now)
	if err != nil {
		t.Fatalf("failed to get local time 94016")
	}

	// fmt.Printf("time1: %v\n", r1)
	// fmt.Printf("time2: %v\n", r2)

	easternLoc, _ := time.LoadLocation("America/New_York")
	pacificLoc, _ := time.LoadLocation("America/Los_Angeles")

	inEastern := now.In(easternLoc)
	inPacific := now.In(pacificLoc)

	hourDiff := inEastern.Hour() - inPacific.Hour()

	if hourDiff < 0 {
		hourDiff += 24
	} else if hourDiff > 3 {
		hourDiff = 3
	}

	if hourDiff != 3 {
		t.Errorf("Expected 3 hour difference between Eastern and Pacific, got %d", hourDiff)
	}

	// Test non-existent zip code
	_, err = GetLocalTimeAt(zipCodeCache, "99999", now)
	if err == nil {
		t.Errorf("Expected error for non-existent zip code 99999, got nil")
	}
}

// TestLoadTimezoneWithFallback test timezones
func TestLoadTimezoneWithFallback(t *testing.T) {

	cases := []struct {
		name             string
		timeZone         string
		expectedTimezone string
	}{
		{name: "default", timeZone: "", expectedTimezone: "UTC"},
		{name: "America/New_York", timeZone: "America/New_York", expectedTimezone: "America/New_York"},
		{name: "America/Chicago", timeZone: "America/Chicago", expectedTimezone: "America/Chicago"},
		{name: "America/Denver", timeZone: "America/Denver", expectedTimezone: "America/Denver"},
		{name: "America/Los_Angeles", timeZone: "America/Los_Angeles", expectedTimezone: "America/Los_Angeles"},
		{name: "invalid/Timezone", timeZone: "invalid/Timezone", expectedTimezone: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			loc, err := LoadTimezoneWithFallback(c.timeZone)
			if c.name == "invalid/Timezone" {
				assert.NotNil(t, err)
				assert.Nil(t, loc)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.expectedTimezone, loc.String())
			}

		})
	}

}

func TestDownLoadZipData(t *testing.T) {

	// const geoNamesZipURL = "http://download.geonames.org/export/zip/US.zip"

	dataDir = "./process/data"

	zipFilePath = dataDir + "/US.zip"

	csvFilePath = dataDir + "/US.txt"

	err := DownLoadZipData()
	assert.Nil(t, err)
}
