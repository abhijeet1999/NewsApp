package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/abhijeet1999/NewsApp/pkg/controllers"
	"github.com/abhijeet1999/NewsApp/pkg/models"
	"github.com/abhijeet1999/NewsApp/pkg/routes"
	"github.com/abhijeet1999/NewsApp/pkg/utils"
)

func main() {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		panic("News api key is not set")
	}

	fmt.Println("*****************************************************************")
	fmt.Println("Welcome To News App")
	fmt.Println("*****************************************************************")

	// Ensure input.txt exists
	file, ferr := os.Open("input.txt")
	if ferr != nil {
		panic(ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)           // Limit concurrency to 10
	results := make(chan models.Result, 100) // Collect results

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, ",")
		if len(items) < 3 {
			fmt.Printf("Skipping invalid line (needs 3 values): %s\n", line)
			continue
		}

		topic := strings.TrimSpace(items[0])
		days, _ := strconv.Atoi(strings.TrimSpace(items[1]))
		count, _ := strconv.Atoi(strings.TrimSpace(items[2]))

		wg.Add(1)
		sem <- struct{}{} // acquire slot

		go func(topic string, days, count int) {
			defer wg.Done()
			defer func() { <-sem }()

			// Get per-topic mutex to prevent concurrent processing of same topic
			topicMutex := models.GetTopicMutex(topic)
			topicMutex.Lock()
			defer topicMutex.Unlock()

			existing, _ := models.GetNewsBySearchKey(topic)
			length := len(existing.Articles)

			// Ensure SearchKey is always set (fixes "unknown.txt")
			existing.SearchKey = topic

			// If DB has enough
			if length >= count {
				results <- models.Result{Topic: topic, Data: existing, Count: count, Source: "Database"}
				return
			}

			// Fetch missing articles from API
			newArticles := routes.GetNewsapi(topic, count, days, apiKey)

			// Use UPSERT approach - no need for manual deduplication
			// The database UNIQUE constraints will handle duplicates
			controllers.SaveNews(topic, newArticles.Articles)

			// Reload data after UPSERT to get accurate count
			existing, _ = models.GetNewsBySearchKey(topic)
			existing.SearchKey = topic

			// Determine final count to print
			finalCount := len(existing.Articles)
			if finalCount > count {
				finalCount = count
			}

			var source string
			if length == 0 && len(newArticles.Articles) == 0 {
				source = "No Data"
			} else if length == 0 {
				source = "API"
			} else if len(newArticles.Articles) == 0 {
				source = "Database"
			} else {
				source = "Combined"
			}

			results <- models.Result{Topic: topic, Data: existing, Count: finalCount, Source: source}

		}(topic, days, count)

	}

	// Close results after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Ensure /output exists
	outputDir := "output"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Write results to files
	for res := range results {
		// Fallback if still empty
		if strings.TrimSpace(res.Data.SearchKey) == "" {
			res.Data.SearchKey = "unknown"
		}
		utils.PrintArticles(*res.Data, res.Count, res.Source)
	}

	fmt.Println("All tasks complete. Results are stored in", filepath.Join(".", outputDir))
}
