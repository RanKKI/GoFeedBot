package feed

import (
    "GoTeleFeed/config"
    "errors"
    "fmt"
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
    fetcher *Fetcher
    pusher  *Pusher
}

func (f *Feed) Init(config config.Config, bot *tgbotapi.BotAPI) {
    f.fp = gofeed.NewParser()
    f.fp.Client = config.Client
    f.config = config
    f.fetcher = &Fetcher{
        PushChannel:  make(chan *UserFeed),
        FetchChannel: make(chan *Item),
        Config:       config,
    }
    f.pusher = &Pusher{
        Bot:    bot,
        Config: config,
        Stream: f.fetcher.PushChannel,
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

func (f *Feed) StartSchedule() {
    var wg sync.WaitGroup
    wg.Add(1)
    cr := cron.New()

    interval := fmt.Sprintf("@every %s", f.config.Interval)
    log.Printf("Interval %s", interval)
    _, err := cr.AddFunc(interval, f.fetcher.Fetch)
    if err != nil {
        log.Panicln(err)
    }
    cr.Run()
    wg.Wait()
}

func (f *Feed) StartService() {
    if f.fp == nil || f.fetcher == nil || f.pusher == nil {
        panic("You must init feed instance first")
    }
    go f.StartSchedule()
    go f.fetcher.StartFetchServices()
    go f.pusher.StartPushServices()
}
