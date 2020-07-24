package main

import (
    "GoTeleFeed/command"
    "GoTeleFeed/database"
    "GoTeleFeed/feed"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/robfig/cron/v3"
    "log"
    "sync"

    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

func newBot(config *Config) *tgbotapi.BotAPI {
    log.Println("Starting Bot...")

    bot, err := tgbotapi.NewBotAPIWithClient(config.Token, config.Client)
    if err != nil {
        panic(err)
    }

    bot.Debug = config.Debug

    log.Printf("Authorized on account %s", bot.Self.UserName)

    return bot
}

func startListening(bot *tgbotapi.BotAPI) {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, err := bot.GetUpdatesChan(u)
    if err != nil {
        panic(err)
    }

    log.Println("Start listening")

    for update := range updates {
        go command.HandleCommand(bot, &update)
    }
}

func startFeedServices(bot *tgbotapi.BotAPI) {
    stream := feed.PushFeedServices(bot)

    job := func() {
        log.Println("--------------------------------------")
        feed.CheckForUpdates(stream)
    }

    go func(ch chan *feed.UserFeeds) {
        wg := sync.WaitGroup{}
        wg.Add(1)
        cr := cron.New()
        // run every 20 minutes
        _, err := cr.AddFunc("*/20 * * * *", job)
        if err != nil {
            log.Panicln(err)
        }
        cr.Run()
        wg.Wait()
    }(stream)
}

func main() {
    // config
    config := loadConfig()
    database.Init(config.Debug)
    bot := newBot(config)
    feed.InitFeedParser(config.Client)

    // start services
    startFeedServices(bot)
    startListening(bot)
}
