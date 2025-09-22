package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/abhijeet1999/NewsApp/pkg/controllers"
	"github.com/abhijeet1999/NewsApp/pkg/models"
)

func GetNewsapi(searchTopic string, count int, daysAgo int, apiKey string) models.NewsData {
	fmt.Println()
	fmt.Println("Fetching articles count:", count)
	baseURL := "https://newsapi.org/v2/everything"

	// Calculate the date n days ago
	fromDate := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

	// Build query params safely
	params := url.Values{}
	params.Add("q", searchTopic)
	params.Add("from", fromDate)
	params.Add("sortBy", "popularity")
	params.Add("apiKey", apiKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error fetching news:", err)
		return models.NewsData{SearchKey: searchTopic, Articles: []models.Article{}}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return models.NewsData{SearchKey: searchTopic, Articles: []models.Article{}}
	}

	var apiResp struct {
		Articles []models.Article `json:"articles"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return models.NewsData{SearchKey: searchTopic, Articles: []models.Article{}}
	}

	// Limit articles safely: if fewer than count, just return whatever is available
	limitedArticles := apiResp.Articles
	if len(apiResp.Articles) > count {
		limitedArticles = apiResp.Articles[:count]
	}

	news := models.NewsData{
		SearchKey: searchTopic,
		Articles:  limitedArticles,
	}

	fmt.Printf("Result from Internet: %d articles\n", len(limitedArticles))
	controllers.SaveNews(searchTopic, limitedArticles) // save whatever was fetched

	return news
}
