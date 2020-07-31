package feed

import (
    "GoTeleFeed/config"
    "errors"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/mmcdole/gofeed"
    "github.com/robfig/cron/v3"
    "log"
    "net/url"
    "sync"
)

var Instance Feed

type Feed struct {
    fp      *gofeed.Parser
    config  config.Config
    ch      chan *UserFeed
    fetcher *Fetcher
}

func (f *Feed) Init(config config.Config, bot *tgbotapi.BotAPI) {
    f.fp = gofeed.NewParser()
    f.fp.Client = config.Client
    f.config = config
    f.ch = PushFeedServices(bot)
    f.fetcher = &Fetcher{
        PushChannel:  f.ch,
        FetchChannel: make(chan *Items),
    }
}

func (f *Feed) TestURL(feedURL string) (*gofeed.Feed, error) {
    log.Printf("Testing %s", feedURL)
    if f.fp == nil {
        return nil, errors.New("you should init feed parser first")
    }
    _, err := url.ParseRequestURI(feedURL)
    if err != nil {
        return nil, err
    }
    feed, err := f.fp.ParseURL(feedURL)
    if feed != nil && feed.Items[0].Published == "" {
        return feed, errors.New("this feed doesn't have publish attribute")
    }
    return feed, err
}

func (f *Feed) StartService() {
    go f.fetcher.StartFetchServices()

    job := func() {
        log.Println("--------------------------------------")
        f.fetcher.Fetch()
    }

    go func(ch chan *UserFeed) {
        wg := sync.WaitGroup{}
        wg.Add(1)
        cr := cron.New()
        // run every 20 minutes
        _, err := cr.AddFunc("*/5 * * * *", job)
        if err != nil {
            log.Panicln(err)
        }
        cr.Run()
        wg.Wait()
    }(f.ch)
}
