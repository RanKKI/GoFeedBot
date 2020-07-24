package database

import (
    "errors"
    "github.com/mmcdole/gofeed"
    "time"
)

func AddFeed(chatID int64, feed *gofeed.Feed, url string) error {
    f := &Feed{
        URL:    url,
        ChatID: chatID,
    }

    db.Where(f).First(f)

    if f.ID != 0 {
        return errors.New("already exists")
    }

    f.Title = feed.Title
    if feed.Author != nil {
        f.AuthorName = feed.Author.Name
    }
    f.LastCheckAt = time.Now().UTC()
    db.Create(f)
    return nil
}

func QueryFeeds(chatID int64) []*Feed {
    var feeds []*Feed
    db.Where(&Feed{ChatID: chatID}).Find(&feeds)
    return feeds
}

func QueryAllFeeds() []*Feed {
    var feeds []*Feed
    db.Find(&feeds)
    return feeds
}

func QueryAllLinks() []string {
    var links []string
    db.Table(TableName).Select("distinct url").Pluck("url", &links)
    return links
}

func GetUpdateTime(chatID int64) time.Time {
    feed := Feed{
        ChatID: chatID,
    }
    db.Where(&feed).First(&feed)
    return feed.LastCheckAt.UTC()
}

func UpdateTime(chatID int64) {
    db.Table(TableName).Where(&Feed{
        ChatID: chatID,
    }).Update(Feed{
        LastCheckAt: time.Now(),
    })
}
