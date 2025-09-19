package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/abhijeet1999/NewsApp/pkg/models"
	"github.com/abhijeet1999/NewsApp/pkg/routes"
	"github.com/abhijeet1999/NewsApp/pkg/utils"
)

func main() {
	fmt.Println("*****************************************************************")
	fmt.Println("Welcome To News App")
	fmt.Println("*****************************************************************")

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

			existing, _ := models.GetNewsBySearchKey(topic)
			length := len(existing.Articles)

			// If DB has enough
			if length >= count {
				results <- models.Result{Topic: topic, Data: existing, Count: count, Source: "Database"}
				return
			}

			// Fetch missing articles from API
			missing := count

			newArticles := routes.GetNewsapi(topic, missing, days)

			// Merge unique articles from API
			urlMap := make(map[string]bool)
			for _, a := range existing.Articles {
				urlMap[a.URL] = true
			}

			for _, a := range newArticles.Articles {
				if !urlMap[a.URL] {
					existing.Articles = append(existing.Articles, a)
				}
			}

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

	// Write results to files in /output/
	for res := range results {
		utils.PrintArticles(*res.Data, res.Count, res.Source)
	}

	fmt.Println("All tasks complete.")
}
