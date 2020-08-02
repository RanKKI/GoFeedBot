package feed

import (
    "GoTeleFeed/config"
    "GoTeleFeed/model"
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
    URL   string
    Items []*gofeed.Item
}

type Fetcher struct {
    PushChannel  chan *UserFeed
    FetchChannel chan *Item
    Config       config.Config
    Subscribes   map[string]map[int64]bool
}

func (f *Fetcher) getLatestTime(t1 *time.Time, t2 *time.Time) time.Time {
    t3 := t1
    if t1 != nil && t2 != nil && t2.After(*t1) {
        t3 = t2
    } else if t3 == nil && t2 != nil {
        t3 = t2
    } else if t3 == nil {
        log.Panicf("error on parsing time, where t1=%s, t2=%s", t1, t2)
    }
    return *t3
}

func (f *Fetcher) fetchURL(url string, wg *sync.WaitGroup) {
    defer wg.Done()
    if f.Config.Debug {
        log.Printf("Checking %s", url)
    }
    feed, err := Instance.fp.ParseURL(url)

    if err != nil {
        log.Panicln(err)
    }

    lastUpdatedTime := f.getLatestTime(feed.UpdatedParsed, feed.PublishedParsed)
    if f.Config.Debug {
        log.Printf("Totoal %d items of %s, Updated at %s", len(feed.Items), feed.Title, lastUpdatedTime)
    }

    f.FetchChannel <- &Item{
        URL:   url,
        Items: feed.Items,
    }
}

func (f *Fetcher) Fetch() {
    // Update subscribes
    for _, feed := range model.QueryFeeds() {
        if f.Subscribes[feed.URL] == nil {
            f.Subscribes[feed.URL] = make(map[int64]bool)
        }
        f.Subscribes[feed.URL][feed.ChatID] = true
    }

    log.Printf("Requesting %d links", len(f.Subscribes))

    var wg sync.WaitGroup
    for _, url := range model.QueryFeedURLs() {
        wg.Add(1)
        go f.fetchURL(url, &wg)
    }
    wg.Wait()
}

func (f *Fetcher) StartFetchServices() {
    for item := range f.FetchChannel {
        for chatID := range f.Subscribes[item.URL] {
            go f.FilterAndSend(chatID, item)
        }
    }

}

func (f *Fetcher) FilterAndSend(chatID int64, item *Item) {
    lastCheckTime := model.QueryFeed(model.Feed{ChatID: chatID}).LastCheckAt

    for _, item := range item.Items {

        if item.PublishedParsed == nil {
            log.Printf("Item %s does not have published time", item.Title)
        } else if item.PublishedParsed.After(lastCheckTime) {
            f.PushChannel <- &UserFeed{
                ChatID: chatID,
                Item:   item,
            }
        } else {
            // since the items are in-order
            // if one of the items is published before the `lastCheckTime`
            // ignored all of rest items
            break
        }
    }

    model.UpdateTime(chatID, item.URL)
}
