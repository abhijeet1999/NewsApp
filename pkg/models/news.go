package models

import (
	"fmt"

	"github.com/abhijeet1999/NewsApp/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Article struct {
	gorm.Model
	NewsDataID  uint   `gorm:"index"`
	Title       string `json:"title" gorm:"type:varchar(255)"`
	Author      string `json:"author" gorm:"type:varchar(255)"`
	URL         string `json:"url" gorm:"type:text"`
	Description string `json:"description" gorm:"type:text"`
	PublishedAt string `json:"publishedAt" gorm:"type:varchar(100)"`
}

type NewsData struct {
	gorm.Model
	SearchKey string    `gorm:"column:searchkey" json:"searchkey"`
	Articles  []Article `gorm:"foreignKey:NewsDataID"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&NewsData{}, &Article{}) // ✅ migrate both
}

func GetNewsBySearchKey(searchKey string) (*NewsData, *gorm.DB) {
	var getNews NewsData
	result := db.Preload("Articles").Where("searchkey = ?", searchKey).First(&getNews)

	return &getNews, result
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
