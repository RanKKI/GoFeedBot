package feed

import (
    "GoFeedBot/config"
    "GoFeedBot/model"
    "errors"
    "github.com/mmcdole/gofeed"
    "log"
    "sync"
    "time"
)

type UserFeed struct {
    ChatID int64
    Item   *gofeed.Item
}

type Item struct {
    Feed  *model.Feed
    Items []*gofeed.Item
}

type Fetcher struct {
    PushChannel  chan *UserFeed
    FetchChannel chan *Item
    Config       config.Config
}

func (fetcher *Fetcher) getLatestTime(t1 *time.Time, t2 *time.Time) (*time.Time, error) {
    t3 := t1
    if t1 != nil && t2 != nil && t2.After(*t1) {
        t3 = t2
    } else if t3 == nil && t2 != nil {
        t3 = t2
    } else if t3 == nil {
        return nil, errors.New("both given time are nil")
    }
    return t3, nil
}

func (fetcher *Fetcher) fetchURL(feed *model.Feed, wg *sync.WaitGroup) {
    defer wg.Done()
    url := feed.URL
    if fetcher.Config.Debug {
        log.Printf("Checking %s", url)
    }
    f, err := Instance.fp.ParseURL(url)

    if err != nil {
        log.Printf("error on fetching %s, %s", feed.Title, err.Error())
        return
    }

    lastUpdatedTime, err := fetcher.getLatestTime(f.UpdatedParsed, f.PublishedParsed)
    if err != nil {
        log.Printf("error on fetching %s, %s", feed.Title, err.Error())
        return
    }
    if fetcher.Config.Debug {
        log.Printf("Totoal %d items of %s, Updated at %s", len(f.Items), feed.Title, lastUpdatedTime)
    }

    fetcher.FetchChannel <- &Item{
        Feed:  feed,
        Items: f.Items,
    }
}

func (fetcher *Fetcher) Fetch() {
    // Update subscribes
    feeds := model.GetAllFeeds()

    log.Printf("Requesting %d links", len(feeds))

    var wg sync.WaitGroup
    for _, feed := range feeds {
        // if there no subscriber of the feed, skip
        if len(feed.Subscribers) == 0 {
            continue
        }
        wg.Add(1)
        go fetcher.fetchURL(feed, &wg)
    }
    wg.Wait()
}

func (fetcher *Fetcher) StartFetchServices() {
    for item := range fetcher.FetchChannel {
        for _, sub := range item.Feed.Subscribers {
            go fetcher.FilterAndSend(sub.ChatID, item)
        }
    }

}

func (fetcher *Fetcher) FilterAndSend(chatID int64, item *Item) {
    lastCheckTime := item.Feed.LastCheckAt

    for _, item := range item.Items {

        if item.PublishedParsed == nil {
            log.Printf("Item %s does not have published time", item.Title)
        } else if item.PublishedParsed.After(lastCheckTime) && !model.ItemSentBefore(item) {
            fetcher.PushChannel <- &UserFeed{
                ChatID: chatID,
                Item:   item,
            }
            model.NewItem(item)
        } else {
            // since the items are in-order
            // if one of the items is published before the `lastCheckTime`
            // ignored all of rest items
            break
        }
    }

    item.Feed.UpdateCheckTime()
}
