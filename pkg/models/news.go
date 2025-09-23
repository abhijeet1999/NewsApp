package models

import (
	"fmt"
	"sync"

	"github.com/abhijeet1999/NewsApp/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// Per-topic mutex map to prevent concurrent processing of same topic
var topicMutexes = make(map[string]*sync.Mutex)
var mutexMapMutex sync.RWMutex

type Article struct {
	gorm.Model
	NewsDataID  uint   `gorm:"index"`
	Title       string `json:"title" gorm:"type:varchar(255)"`
	Author      string `json:"author" gorm:"type:varchar(255)"`
	URL         string `json:"url" gorm:"type:text;unique_index"` // UNIQUE index on URL
	Description string `json:"description" gorm:"type:text"`
	PublishedAt string `json:"publishedAt" gorm:"type:varchar(100)"`
}

type NewsData struct {
	gorm.Model
	SearchKey string    `gorm:"column:searchkey;unique_index" json:"searchkey"` // UNIQUE index on SearchKey
	Articles  []Article `gorm:"foreignKey:NewsDataID"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&NewsData{}, &Article{}) // ✅ migrate both
}

// GetTopicMutex returns a mutex for the given topic to prevent concurrent processing
func GetTopicMutex(topic string) *sync.Mutex {
	mutexMapMutex.Lock()
	defer mutexMapMutex.Unlock()

	if mutex, exists := topicMutexes[topic]; exists {
		return mutex
	}

	// Create new mutex for this topic
	mutex := &sync.Mutex{}
	topicMutexes[topic] = mutex
	return mutex
}

func GetNewsBySearchKey(searchKey string) (*NewsData, *gorm.DB) {
	var getNews NewsData
	result := db.Preload("Articles").Where("searchkey = ?", searchKey).First(&getNews)

	return &getNews, result
}

// UpsertNewsData creates or updates NewsData with UPSERT logic
func UpsertNewsData(searchKey string) (*NewsData, error) {
	var newsData NewsData

	// Try to find existing record
	result := db.Where("searchkey = ?", searchKey).First(&newsData)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new record
			newsData = NewsData{SearchKey: searchKey}
			if err := db.Create(&newsData).Error; err != nil {
				return nil, err
			}
			fmt.Printf("Created new topic: %s\n", searchKey)
		} else {
			return nil, result.Error
		}
	} else {
		fmt.Printf("Found existing topic: %s\n", searchKey)
	}

	return &newsData, nil
}

// UpsertArticle creates or updates Article with UPSERT logic using UNIQUE index
func UpsertArticle(article Article) error {
	// Use GORM's FirstOrCreate with UNIQUE constraint
	result := db.Where("url = ?", article.URL).FirstOrCreate(&article)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		fmt.Printf("Created new article: %s\n", article.Title)
	} else {
		fmt.Printf("Article already exists: %s\n", article.Title)
	}

	return nil
}

// SaveArticlesWithUpsert saves articles using UPSERT to prevent duplicates
func SaveArticlesWithUpsert(searchKey string, articles []Article) error {
	// Get or create NewsData
	newsData, err := UpsertNewsData(searchKey)
	if err != nil {
		return err
	}

	// Save each article with UPSERT
	var savedCount int
	for _, article := range articles {
		article.NewsDataID = newsData.ID
		if err := UpsertArticle(article); err != nil {
			fmt.Printf("Error saving article %s: %v\n", article.Title, err)
			continue
		}
		savedCount++
	}

	fmt.Printf("Saved %d articles for topic: %s\n", savedCount, searchKey)
	return nil
}

func (b *NewsData) CreateNews() {
	db.Create(&b)
	fmt.Println("Saved in Database")
}

func DB() *gorm.DB {
	return db
}

type Result struct {
	Topic  string
	Data   *NewsData
	Count  int
	Source string
}
