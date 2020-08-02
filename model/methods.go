package model

import (
    "errors"
    "time"
)

func AddSub(chatID int64, url, title, author string) error {
    f := Feed{
        URL:    url,
        ChatID: chatID,
    }

    db.Where(&f).First(&f)

    if f.ID != 0 {
        return errors.New("already exists")
    }

    f.Title = title
    f.AuthorName = author
    f.LastCheckAt = time.Now().UTC()
    db.Create(&f)
    return nil
}

func DeleteSub(chatID int64, id int64) error {
    feed := Feed{
        ID: id,
    }
    db.Where(&feed).First(&feed)
    if feed.ChatID != chatID {
        return errors.New("invalid ID")
    }
    db.Delete(&feed)
    return nil
}

func QueryFeeds(filters ...Feed) []*Feed {
    var feeds []*Feed
    query := db
    for _, filter := range filters {
        query = query.Where(filter)
    }
    query.Find(&feeds)
    return feeds
}

func QueryFeed(filters ...interface{}) *Feed {
    feed := Feed{}
    query := db
    for _, filter := range filters {
        query = query.Where(filter)
    }
    query.First(&feed)
    return &feed
}

func QueryFeedURLs() []string {
    var links []string
    db.Table(TableName).Select("distinct url").Pluck("url", &links)
    return links
}

func UpdateTime(chatID int64, url string) {
    db.Table(TableName).Where(&Feed{
        ChatID: chatID,
        URL:    url,
    }).Update(Feed{
        LastCheckAt: time.Now(),
    })
}
