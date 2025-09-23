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

// SaveNewsWithUpsert uses UPSERT approach with UNIQUE constraints to prevent duplicates
func SaveNewsWithUpsert(searchKey string, articles []models.Article) {
	// Use the new UPSERT function that handles UNIQUE constraints
	err := models.SaveArticlesWithUpsert(searchKey, articles)
	if err != nil {
		fmt.Printf("Error saving articles for %s: %v\n", searchKey, err)
	}
}

// SaveNews is the legacy function - now uses UPSERT approach
func SaveNews(searchKey string, articles []models.Article) {
	SaveNewsWithUpsert(searchKey, articles)
}
