package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abhijeet1999/NewsApp/pkg/models"
)

func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}

	}
}
func SplitTopics(topics string) []string {
	topicList := strings.Split(topics, ",")
	for i := range topicList {
		topicList[i] = strings.TrimSpace(topicList[i])
	}
	return topicList
}

func PrintArticles(newsdata models.NewsData, count int, source string) {
	// Ensure /output directory exists
	outputDir := "output"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Create file path (e.g., /output/apple.txt)
	filename := strings.ReplaceAll(strings.ToLower(newsdata.SearchKey), " ", "_") + ".txt"
	filepath := filepath.Join(outputDir, filename)

	// Read any existing content so we can prepend the latest results
	var existingContent []byte
	if _, statErr := os.Stat(filepath); statErr == nil {
		if content, readErr := os.ReadFile(filepath); readErr == nil {
			existingContent = content
		}
	}

	// Recreate the file and write the latest block first
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filepath, err)
		return
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Determine how many articles will actually be printed
	displayedCount := count
	if len(newsdata.Articles) < displayedCount {
		displayedCount = len(newsdata.Articles)
	}

	// Write header with source
	fmt.Fprintf(f, "**************************************************************\n")
	fmt.Fprintf(f, "Timestamp: %s\n", timestamp)
	fmt.Fprintf(f, "Results for: %-20s | Showing %d articles | Source: %s\n", newsdata.SearchKey, displayedCount, source)
	fmt.Fprintf(f, "**************************************************************\n\n")

	// Write each article
	for i, article := range newsdata.Articles {
		if i >= displayedCount {
			break
		}
		fmt.Fprintf(f, "[%d] Title       : %s\n", i+1, article.Title)
		fmt.Fprintf(f, "    Author      : %s\n", article.Author)
		fmt.Fprintf(f, "    Description : %s\n", article.Description)
		fmt.Fprintf(f, "    URL         : %s\n\n", article.URL)
	}

	// Footer
	fmt.Fprintf(f, "**************************** END ****************************\n\n")

	// If there was existing content, append it below with a separator
	if len(existingContent) > 0 {
		fmt.Fprintf(f, "--- PREVIOUS RESULTS ---\n\n")
		f.Write(existingContent)
	}
}
