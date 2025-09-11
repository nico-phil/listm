package tz

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// LoadZipCodeData loads zip code data from the GeoNames database
func LoadZipCodeData() (*ZipCodeCache, error) {
	// Create the data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Check if we need to download the data
	shouldDownload, err := shouldDownloadZipData()
	if err != nil {
		log.Printf("Error checking if zip data should be downloaded: %v", err)
		log.Println("Will attempt to use existing data if available")
	}

	client := http.Client{}
	if shouldDownload {
		log.Println("Downloading zip code data...")
		if err := downloadZipData(context.Background(), client, geoNamesZipURL, zipFilePath); err != nil {
			return nil, fmt.Errorf("failed to download zip data: %v", err)
		}

		log.Println("Extracting zip code data...")
		if err := extractZipData(); err != nil {
			return nil, fmt.Errorf("failed to extract zip data: %v", err)
		}
	} else {
		log.Println("Using existing zip code data")
	}

	// Load the data into memory
	return loadZipCodeDataFromCSV()

}

// shouldDownloadZipData checks if we need to download the zip data
func shouldDownloadZipData() (bool, error) {
	// If the CSV file doesn't exist, we definitely need to download
	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		return true, nil
	}

	// Check the file size on the server
	resp, err := http.Head(geoNamesZipURL)
	if err != nil {
		return false, fmt.Errorf("failed to check remote file size: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-OK status: %s", resp.Status)
	}

	contentLength := resp.Header.Get("Content-Length")
	remoteSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse content length: %v", err)
	}

	// Get the local file size
	fileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to stat local file: %v", err)
	}

	localSize := fileInfo.Size()

	// If the file sizes are different, we need to download
	return remoteSize != localSize, nil
}

// loadZipCodeDataFromCSV loads the zip code data from the CSV file
func loadZipCodeDataFromCSV() (*ZipCodeCache, error) {
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	cache := NewZipCodeCache()
	scanner := bufio.NewScanner(file)

	// The file format is tab-separated, with the following columns:
	// 0: country code
	// 1: postal code
	// 2: place name
	// 3: admin name1 (state)
	// 4: admin code1 (state code)
	// 5: admin name2 (county)
	// 6: admin code2
	// 7: admin name3
	// 8: admin code3
	// 9: latitude
	// 10: longitude
	// 11: accuracy
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		if len(fields) < 12 {
			continue
		}

		// Skip the first line if it's a header
		if fields[0] == "country code" {
			continue
		}

		zipCode := fields[1]
		lat, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			continue
		}

		long, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			continue
		}

		// Determine the time zone based on the longitude
		// This is a simple approximation - in a production system you would want to use a
		// proper timezone database that maps coordinates to timezone IDs
		timeZone := determineTimeZone(long)

		info := &ZipCodeInfo{
			ZipCode:   zipCode,
			Latitude:  lat,
			Longitude: long,
			City:      fields[2],
			State:     fields[3],
			TimeZone:  timeZone,
		}

		cache.Set(zipCode, info)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning CSV file: %v", err)
	}

	log.Printf("Loaded information for %d zip codes", len(cache.cache))
	return cache, nil
}

// determineTimeZone determines the time zone based on longitude
// This is a very basic approximation
func determineTimeZone(longitude float64) string {
	// US longitude ranges from about -125 (West Coast) to -67 (East Coast)
	// Time zones are roughly:
	// Eastern: -82 to -67
	// Central: -100 to -82
	// Mountain: -115 to -100
	// Pacific: -125 to -115

	if longitude >= -82 {
		return "America/New_York" // Eastern
	} else if longitude >= -100 {
		return "America/Chicago" // Central
	} else if longitude >= -115 {
		return "America/Denver" // Mountain
	} else {
		return "America/Los_Angeles" // Pacific
	}
}

// extractZipData extracts the zip code data from the downloaded zip file
func extractZipData() error {
	// Open the zip file
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipFile.Close()

	// Find the US.txt file
	var usFile *zip.File
	for _, file := range zipFile.File {
		if file.Name == "US.txt" {
			usFile = file
			break
		}
	}

	if usFile == nil {
		return fmt.Errorf("US.txt not found in zip file")
	}

	// Open the US.txt file inside the zip
	usFileReader, err := usFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open US.txt inside zip: %v", err)
	}
	defer usFileReader.Close()

	// Create the output file
	outFile, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Copy the file to the output
	_, err = io.Copy(outFile, usFileReader)
	return err
}
