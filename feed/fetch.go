package feed

import (
    "GoTeleFeed/database"
    "github.com/mmcdole/gofeed"
    "log"
    "sync"
)

type UserFeeds struct {
    ChatID int64
    Items  []*gofeed.Item
}

func addItems(w *sync.WaitGroup, lock *sync.Mutex, link string, m *map[string][]*gofeed.Item) {
    log.Printf("Checking %s", link)
    feed, err := fp.ParseURL(link)
    if err != nil {
        log.Println(err)
        return
    }

    // add items to map
    lock.Lock()
    (*m)[link] = feed.Items
    lock.Unlock()

    lastUpdatedTime := feed.UpdatedParsed
    if feed.UpdatedParsed != nil && feed.PublishedParsed != nil {
        if feed.PublishedParsed.Sub(*feed.UpdatedParsed) > 0 {
            lastUpdatedTime = feed.PublishedParsed
        }
    } else if lastUpdatedTime == nil && feed.PublishedParsed != nil {
        lastUpdatedTime = feed.PublishedParsed
    }

    if lastUpdatedTime != nil {
        log.Printf("Totoal %d items of %s, Updated at %s", len(feed.Items), feed.Title, lastUpdatedTime)
    } else {
        log.Printf("Totoal %d items of %s", len(feed.Items), feed.Title)
    }
    w.Done()
}

func fetchAllUserSubscribe() *map[int64][]string {
    m := map[int64][]string{}
    for _, f := range database.QueryAllFeeds() {
        m[f.ChatID] = append(m[f.ChatID], f.URL)
    }
    return &m
}

func fetchAllSubscribeItems() *map[string][]*gofeed.Item {
    m := map[string][]*gofeed.Item{}
    wg := sync.WaitGroup{}
    lock := sync.Mutex{}
    for _, link := range database.QueryAllLinks() {
        wg.Add(1)
        go addItems(&wg, &lock, link, &m)
    }
    wg.Wait()
    return &m
}

func CheckUpdates() []*UserFeeds {
    userFeeds := *fetchAllUserSubscribe()
    feedItems := *fetchAllSubscribeItems()

    // filter data
    var ret []*UserFeeds

    for chatID, links := range userFeeds {
        userFeeds := UserFeeds{
            ChatID: chatID,
            Items:  []*gofeed.Item{},
        }
        lastCheckTime := database.GetUpdateTime(chatID)
        log.Printf("User %d last check at %s", chatID, lastCheckTime)
        for _, link := range links {
            for _, item := range feedItems[link] {
                if item.PublishedParsed == nil {
                    log.Printf("Item %s does not have published time", item.Title)
                } else if item.PublishedParsed.Sub(lastCheckTime) >= 0 {
                    userFeeds.Items = append(userFeeds.Items, item)
                }
            }
        }
        database.UpdateTime(chatID)
        ret = append(ret, &userFeeds)
        log.Printf("%d items have to sent to %d", len(userFeeds.Items), userFeeds.ChatID)
    }
    return ret
}
