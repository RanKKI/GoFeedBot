package feed

import (
    "GoTeleFeed/database"
    "github.com/mmcdole/gofeed"
    "log"
    "sync"
)

type UserFeed struct {
    ChatID int64
    Item   *gofeed.Item
}

type Items struct {
    URL   string
    Items []*gofeed.Item
}

type Fetcher struct {
    PushChannel  chan *UserFeed
    FetchChannel chan *Items
    Debug        bool
}

func (f *Fetcher) fetchSubscribes() map[string][]int64 {
    urlUsers := map[string][]int64{}

    for _, feed := range database.QueryAllFeeds() {
        urlUsers[feed.URL] = append(urlUsers[feed.URL], feed.ChatID)
    }

    return urlUsers
}

func (f *Fetcher) fetchURL(url string, wg *sync.WaitGroup) {
    defer wg.Done()
    log.Printf("Checking %s", url)
    feed, err := Instance.fp.ParseURL(url)

    if err != nil {
        log.Panicln(err)
    }

    lastUpdatedTime := feed.UpdatedParsed
    if feed.UpdatedParsed != nil && feed.PublishedParsed != nil {
        if feed.PublishedParsed.After(*feed.UpdatedParsed) {
            lastUpdatedTime = feed.PublishedParsed
        }
    } else if lastUpdatedTime == nil && feed.PublishedParsed != nil {
        lastUpdatedTime = feed.PublishedParsed
    }

    if lastUpdatedTime != nil {
        log.Printf("Totoal %d items of %s, Updated at %s", len(feed.Items), feed.Title, lastUpdatedTime)
    } else {
        // Should't happen
        log.Printf("Totoal %d items of %s", len(feed.Items), feed.Title)
    }

    f.FetchChannel <- &Items{
        URL:   url,
        Items: feed.Items,
    }
}

func (f *Fetcher) Fetch() {
    var wg sync.WaitGroup
    for _, url := range database.QueryAllLinks() {
        wg.Add(1)
        go f.fetchURL(url, &wg)
    }
    wg.Wait()
}

func (f *Fetcher) StartFetchServices() {
    urls := f.fetchSubscribes()

    for item := range f.FetchChannel {
        for _, chatID := range urls[item.URL] {
            lastCheckTime := database.GetUpdateTime(chatID)
            hasItem := false
            for _, item := range item.Items {

                if item.PublishedParsed == nil {
                    log.Printf("Item %s does not have published time", item.Title)
                } else if item.PublishedParsed.Sub(lastCheckTime) >= 0 {
                    hasItem = true
                    f.PushChannel <- &UserFeed{
                        ChatID: chatID,
                        Item:   item,
                    }
                }
            }

            if hasItem {
                database.UpdateTime(chatID, item.URL)
            }

        }
    }

}
