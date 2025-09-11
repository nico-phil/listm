package tz

import (
	"fmt"
	"sync"
	"time"
)

var (
	// GeoNamesZipURL is the URL to download zip code data
	geoNamesZipURL = "http://download.geonames.org/export/zip/US.zip"
	// DataDir is the directory where we'll store the downloaded data
	dataDir = "data"
	// ZipFilePath is the path to the downloaded zip file
	zipFilePath = dataDir + "/US.zip"
	// CSVFilePath is the path to the extracted CSV file
	csvFilePath = dataDir + "/US.txt"
)

// ZipCode contains information about zipcode
type ZipCodeInfo struct {
	ZipCode   string
	Latitude  float64
	Longitude float64
	State     string
	City      string
	TimeZone  string
}

// ZipcodeCache maps zipcode to their information
type ZipCodeCache struct {
	cache map[string]*ZipCodeInfo
	mu    sync.RWMutex
}

// NewZipcodeCache stores zipcode informations
func NewZipCodeCache() *ZipCodeCache {
	return &ZipCodeCache{
		cache: make(map[string]*ZipCodeInfo),
	}
}

// GetZipcode return zipcodeinfo for a given zipcode
func (z *ZipCodeCache) getZipcode(zipCode string) (*ZipCodeInfo, bool) {
	z.mu.RLock()
	zipCodeInfo, ok := z.cache[zipCode]
	z.mu.RUnlock()
	return zipCodeInfo, ok
}

// setZipCode adds or updates zipcode information
func (z *ZipCodeCache) setZipCode(zipCode string, zipCodeInfo *ZipCodeInfo) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.cache[zipCode] = zipCodeInfo
}

func (z *ZipCodeCache) Set(zipCode string, zipCodeInfo *ZipCodeInfo) {
	z.setZipCode(zipCode, zipCodeInfo)
}

// func(z *ZipCodeCache) GetLocalTime(zip)

func GetLocalTimeAt(zipCodeCache *ZipCodeCache, zipCode string, utcTime time.Time) (time.Time, error) {
	info, ok := zipCodeCache.cache[zipCode]
	if !ok {
		return time.Time{}, fmt.Errorf("zipcode not found in the cache")
	}

	loc, err := LoadTimezoneWithFallback(info.TimeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone %v", err)
	}

	return utcTime.In(loc), nil
}

func LoadTimezoneWithFallback(timezone string) (*time.Location, error) {
	// Try loading the timezone by name first
	loc, err := time.LoadLocation(timezone)
	if err == nil {
		return loc, nil
	}

	// If that fails, use IANA timezone database-backed time zones
	// These correctly handle DST transitions automatically
	switch timezone {
	case "America/New_York":
		return time.LoadLocation("EST5EDT")
	case "America/Chicago":
		return time.LoadLocation("CST6CDT")
	case "America/Denver":
		return time.LoadLocation("MST7MDT")
	case "America/Los_Angeles":
		return time.LoadLocation("PST8PDT")
	default:
		return nil, fmt.Errorf("unknown time zone: %s", timezone)
	}
}
