package tz

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownLoadZipData downloads the zip code data GeoNames
func DownLoadZipData() error {

	// create the output file
	file, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file %v", err)
	}

	// download the file
	resp, err := http.Get(geoNamesZipURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status: %d", resp.StatusCode)
	}

	//copy the file to the ouput
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
