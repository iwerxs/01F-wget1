package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// downloadFile downloads a file from a URL and saves it locally
func downloadFile(url string, outputPath string) error {
	startTime := time.Now()
	fmt.Printf("Start Time: %s\n", startTime.Format("2006-01-02 15:04:05"))

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	defer resp.Body.Close()

	// Check for HTTP status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	fmt.Printf("Status: %s\n", resp.Status)

	// Get content length
	contentLength := resp.ContentLength
	fmt.Printf("File Size: %.2f MB (%d bytes)\n", float64(contentLength)/(1024*1024), contentLength)

	// Determine file name
	if outputPath == "" {
		outputPath = path.Base(url)
	}

	// Create file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	// Track progress
	buffer := make([]byte, 4096)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := outFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to file: %v", writeErr)
			}
			downloaded += int64(n)

			// Show progress
			progress := float64(downloaded) / float64(contentLength) * 100
			fmt.Printf("\rDownloading... %.2f%% (%d/%d bytes)", progress, downloaded, contentLength)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error during download: %v", err)
		}
	}
	fmt.Println("\nDownload Complete!")

	// Show completion time
	endTime := time.Now()
	fmt.Printf("Completion Time: %s\n", endTime.Format("2006-01-02 15:04:05"))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-wget <URL> [output_path]")
		os.Exit(1)
	}

	url := os.Args[1]
	outputPath := ""
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	err := downloadFile(url, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
