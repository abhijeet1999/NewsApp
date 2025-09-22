package controllers

import (
	"fmt"
	"net/http"

	"github.com/abhijeet1999/NewsApp/pkg/models"
	"github.com/abhijeet1999/NewsApp/pkg/utils"
)

var NewBook models.NewsData

func GetBookById(SeachKey string, count int16) int {

	bookDetails, _ := models.GetNewsBySearchKey(SeachKey)
	fmt.Println("Result of ", SeachKey)
	fmt.Println(bookDetails.Articles)
	if count > int16(len(bookDetails.Articles)) {

	}
	return len(bookDetails.Articles)

}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook := &models.NewsData{}
	utils.ParseBody(r, CreateBook)
	CreateBook.CreateNews()

}

func SaveNews(searchKey string, articles []models.Article) {
	// Check if record already exists
	existing, _ := models.GetNewsBySearchKey(searchKey)

	if existing.ID != 0 {
		// Existing record: append only new articles
		var added int
		for _, article := range articles {
			var found models.Article
			err := models.DB().
				Where("id = ? AND url = ?", existing.ID, article.URL).
				First(&found).Error

			if err != nil { // not found, so create
				article.NewsDataID = existing.ID
				models.DB().Create(&article)
				added++
			}
		}

		fmt.Printf("Added %d new articles for %s (existing record)\n", added, searchKey)

	} else {
		// Create new record
		news := models.NewsData{
			SearchKey: searchKey,
			Articles:  articles,
		}
		news.CreateNews()
		fmt.Printf("Saved %d new articles for %s (new record)\n", len(articles), searchKey)
	}
}
