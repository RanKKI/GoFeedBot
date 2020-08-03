package model

import "errors"

type Subscribe struct {
    ID     int64
    ChatID int64
    FeedID int64
}

func (Subscribe) TableName() string {
    return "subscribes"
}

func (s *Subscribe) Exists() bool {
    db.Where(s).Find(s)
    return s.ID != 0
}

func AddSubscribe(feed *Feed, chatID int64) error {
    subscribe := Subscribe{
        ChatID: chatID,
        FeedID: feed.ID,
    }
    if subscribe.Exists() {
        return errors.New("already exists")
    }
    db.Create(&subscribe)
    return nil
}

func Unsubscribe(chatID, feedID int64) error {
    subscribe := Subscribe{
        ChatID: chatID,
        FeedID: feedID,
    }
    if subscribe.Exists() {
        db.Delete(&subscribe)
        return nil
    }
    return errors.New("not exists")
}

func GetSubscribes(filters ...Subscribe) []*Subscribe {
    subscribes := make([]*Subscribe, 0, 0)
    query := db
    for _, filter := range filters {
        query = query.Where(filter)
    }
    query.Find(&subscribes)
    return subscribes
}

func GetFeedsByChatID(chatID int64) []*Feed {
    feeds := make([]*Feed, 0, 10)
    for _, sub := range GetSubscribes(Subscribe{ChatID: chatID}) {
        feeds = append(feeds, GetFeedByID(sub.FeedID))
    }
    return feeds
}

func GetSubscribesByFeedID(feedID int64) []*Subscribe {
    return GetSubscribes(Subscribe{FeedID: feedID})
}
