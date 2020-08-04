package feed

import (
    "GoFeedBot/config"
    "github.com/stretchr/testify/assert"
    "net/http"
    "testing"
)

func TestFeedData(t *testing.T) {
    ass := assert.New(t)
    var url string
    f := Feed{}
    f.Init(config.Config{}, nil)

    url = "http://baidu.com"

    _, err := f.TestURL(url)
    ass.NotNil(err, "feed parser should be init")

    f.config.Client = &http.Client{}

    _, err = f.TestURL(url)
    ass.NotNilf(err, "%s is not a vaild feed url", url)

    url = "https://feeds.twit.tv/twit.xml"
    feeds, err := f.TestURL(url)
    ass.Nilf(err, "%s should be a vaild url", url)

    if ass.NotNilf(f, "%s should have feed", url) {
        ass.Equal("rss", feeds.FeedType)
        ass.Equal(url, feeds.FeedLink)
    }

}
