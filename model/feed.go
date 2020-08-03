package model

import "time"

type Feed struct {
    ID          int64
    Title       string
    URL         string
    LastCheckAt time.Time
    Subscribers []*Subscribe `gorm:"-"`
}

func (Feed) TableName() string {
    return "feeds"
}

func (f *Feed) Exists() bool {
    db.Where(&f).First(&f)
    return f.ID != 0
}

func (f *Feed) UpdateCheckTime() {
    db.Table(f.TableName()).Where(&Feed{
        ID: f.ID,
    }).Update(Feed{
        LastCheckAt: time.Now().UTC(),
    })
}

func NewFeed(url, title string) *Feed {
    feed := &Feed{URL: url}
    if feed.Exists() {
        return feed
    }
    if title == "" {
        title = url
    }
    feed.Title = title
    feed.LastCheckAt = time.Now().UTC()
    db.Create(feed)
    return feed

}

func GetAllFeeds() []*Feed {
    var feeds []*Feed
    db.Find(&feeds)
    for _, feed := range feeds {
        feed.Subscribers = GetSubscribesByFeedID(feed.ID)
    }
    return feeds
}

func GetFeedByID(id int64) *Feed {
    feed := &Feed{ID: id}
    if feed.Exists() {
        return feed
    }
    return nil
}
