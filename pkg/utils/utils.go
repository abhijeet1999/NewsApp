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

	// Open file in append mode (create if not exists)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filepath, err)
		return
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Write header with source
	fmt.Fprintf(f, "**************************************************************\n")
	fmt.Fprintf(f, "Timestamp: %s\n", timestamp)
	fmt.Fprintf(f, "Results for: %-20s | Showing %d articles | Source: %s\n", newsdata.SearchKey, count, source)
	fmt.Fprintf(f, "**************************************************************\n\n")

	// Write each article
	for i, article := range newsdata.Articles {
		if i >= count {
			break
		}
		fmt.Fprintf(f, "[%d] Title       : %s\n", i+1, article.Title)
		fmt.Fprintf(f, "    Author      : %s\n", article.Author)
		fmt.Fprintf(f, "    Description : %s\n", article.Description)
		fmt.Fprintf(f, "    URL         : %s\n\n", article.URL)
	}

	// Footer
	fmt.Fprintf(f, "**************************** END ****************************\n\n")
}
