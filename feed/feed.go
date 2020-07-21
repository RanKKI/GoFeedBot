package feed

import (
    "errors"
    "github.com/mmcdole/gofeed"
    "log"
    "net/http"
)

var fp *gofeed.Parser

func InitFeedParser(client *http.Client) {
    fp = gofeed.NewParser()
    fp.Client = client
}

func TestFeed(url string) (*gofeed.Feed, error) {
    log.Printf("Testing %s", url)
    feed, err := fp.ParseURL(url)
    if feed != nil && feed.Items[0].Published == "" {
        return feed, errors.New("this feed doesn't have publish attribute")
    }
    return feed, err
}
