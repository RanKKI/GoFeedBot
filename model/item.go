package model

import (
    "crypto/md5"
    "fmt"
    "github.com/mmcdole/gofeed"
)

type Item struct {
    ID   int64
    Hash string `gorm:"type:varchar(32);unique_index"`
}

func HashItem(item *gofeed.Item) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(item.Title+item.Description+item.Link)))
}

func ItemSentBefore(feedItem *gofeed.Item) bool {
    item := Item{
        Hash: HashItem(feedItem),
    }
    db.Where(&item).First(&item)
    return item.ID != 0
}

func NewItem(feedItem *gofeed.Item) {
    item := Item{
        Hash: HashItem(feedItem),
    }
    db.Create(&item)
}
