package feed

import (
    "github.com/stretchr/testify/assert"
    "net/http"
    "testing"
)

func TestFeedData(t *testing.T) {
    ass := assert.New(t)
    var url string

    url = "http://baidu.com"

    _, err := TestFeed(url)
    ass.NotNil(err, "feed parser should be init")

    InitFeedParser(&http.Client{})

    _, err = TestFeed(url)
    ass.NotNilf(err, "%s is not a vaild feed url", url)

    url = "https://feeds.twit.tv/twit.xml"
    f, err := TestFeed(url)
    ass.Nilf(err, "%s should be a vaild url", url)

    if ass.NotNilf(f, "%s should have feed", url) {
        ass.Equal("rss", f.FeedType)
        ass.Equal(url, f.FeedLink)
    }

}
