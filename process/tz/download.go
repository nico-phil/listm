package tz

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPDoer interface {
	Do(*http.Request) *http.Response
}

// DownLoadZipData downloads the zip code data GeoNames
func downloadZipData(ctx context.Context, client http.Client, url, dst string) error {

	// download the file
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK status: %d", resp.StatusCode)
	}

	// create the output file
	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create output file %w", err)
	}
	defer file.Close()

	//copy the file to the ouput
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	return nil
}
