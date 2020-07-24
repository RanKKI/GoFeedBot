package feed

import (
    "errors"
    "github.com/mmcdole/gofeed"
    "log"
    "net/http"
    "net/url"
)

var fp *gofeed.Parser

func InitFeedParser(client *http.Client) {
    fp = gofeed.NewParser()
    fp.Client = client
}

func TestFeed(feedURL string) (*gofeed.Feed, error) {
    log.Printf("Testing %s", feedURL)
    _, err := url.ParseRequestURI(feedURL)
    if err != nil {
        return nil, err
    }
    if fp == nil{
        return nil, errors.New("you should init feed parser first")
    }
    feed, err := fp.ParseURL(feedURL)
    if feed != nil && feed.Items[0].Published == "" {
        return feed, errors.New("this feed doesn't have publish attribute")
    }
    return feed, err
}
