package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vansante/go-ffprobe"
)

func main() {
	var folderPath string
	fmt.Print("Enter the directory with the video files: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	folderPath = scanner.Text()
	var totalDuration float64
	var fileCount int
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		switch strings.ToLower(ext) {
		case ".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm":
		default:
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		data, err := ffprobe.GetProbeDataContext(ctx, path)
		if err != nil {
			log.Printf("Processing error %s: %v", path, err)
			return nil
		}
		duration := data.Format.DurationSeconds
		totalDuration += duration
		fileCount++
		hours := int(duration / 3600)
		minutes := int((duration - float64(hours*3600)) / 60)
		seconds := int(duration) % 60
		fmt.Printf("%s: %02d:%02d:%02d\n", path, hours, minutes, seconds)
		return nil
	})
	if err != nil {
		log.Fatalln("Folder crawl error:", err)
	}
	totalHours := int(totalDuration / 3600)
	totalMinutes := int((totalDuration - float64(totalHours*3600)) / 60)
	totalSeconds := int(totalDuration) % 60
	fmt.Println()
	fmt.Println("Result:")
	fmt.Println("Number of video files:", fileCount)
	fmt.Printf("Total time at normal playback speed: %02d:%02d:%02d\n", totalHours, totalMinutes, totalSeconds)
	totalDuration15x := totalDuration / 1.5
	totalHours = int(totalDuration15x / 3600)
	totalMinutes = int((totalDuration15x - float64(totalHours*3600)) / 60)
	totalSeconds = int(totalDuration15x) % 60
	fmt.Printf("Total time at 1.5x playback speed: %02d:%02d:%02d\n", totalHours, totalMinutes, totalSeconds)
	totalDuration2x := totalDuration / 2
	totalHours = int(totalDuration2x / 3600)
	totalMinutes = int((totalDuration2x - float64(totalHours*3600)) / 60)
	totalSeconds = int(totalDuration2x) % 60
	fmt.Printf("Total time at 2x playback speed: %02d:%02d:%02d\n", totalHours, totalMinutes, totalSeconds)
}
