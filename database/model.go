package database

import (
    "GoTeleFeed/config"
    "github.com/jinzhu/gorm"
    "time"
)

var db *gorm.DB

const TableName = "feeds"

func Init(config config.Config) {
    var err error
    db, err = gorm.Open("sqlite3", "data.db")
    if err != nil {
        panic(err)
    }

    db.AutoMigrate(&Feed{})

    if config.Debug {
        db = db.Debug()
    }
}

type Feed struct {
    ID          int64
    URL         string
    Title       string
    AuthorName  string
    ChatID      int64
    LastCheckAt time.Time
}

func (Feed) TableName() string {
    return TableName
}
