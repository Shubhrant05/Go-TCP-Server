package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

const (
	maxSizeBytes = 5 * 1024 * 1024 // 1KB
)

func compressLogFile(logFilePath string) error {
	logFile, err := os.Open(logFilePath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	compressedLogFilePath := logFilePath + ".gz"
	compressedLogFile, err := os.Create(compressedLogFilePath)
	if err != nil {
		return err
	}
	defer compressedLogFile.Close()

	gzipWriter := gzip.NewWriter(compressedLogFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, logFile)
	if err != nil {
		return err
	}

	return nil
}

func checkAndCompress(logFilePath string) {
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fileSizeBytes := fileInfo.Size()

	if fileSizeBytes > maxSizeBytes {
		err := compressLogFile(logFilePath)
		if err != nil {
			fmt.Println("Error compressing log file:", err)
			return
		}

		err = os.Remove(logFilePath)
		if err != nil {
			fmt.Println("Error removing original log file:", err)
			return
		}

		fmt.Println("Log file compressed successfully.")
	} else {
		fmt.Println("Log file size is below the threshold. No compression needed.")
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run compressLogs.go <logFilePath>")
		os.Exit(1)
	}

	logFilePath := os.Args[1]
	checkAndCompress(logFilePath)
}
