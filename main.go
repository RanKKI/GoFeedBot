package main

import (
    "GoFeedBot/command"
    "GoFeedBot/config"
    "GoFeedBot/feed"
    "GoFeedBot/model"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "log"
)

type FeedBot struct {
    Config *config.Config
    Bot    *tgbotapi.BotAPI
}

func (b *FeedBot) init() {
    if b.Config == nil {
        panic("Config is not provided")
    }
    log.Println("Starting Bot...")
    bot, err := tgbotapi.NewBotAPIWithClient(b.Config.Token, b.Config.Client)
    if err != nil {
        panic(err)
    }
    bot.Debug = b.Config.Debug
    log.Printf("Authorized on account %s", bot.Self.UserName)
    b.Bot = bot
}

func (b *FeedBot) run() {
    if b.Bot == nil {
        panic("You must init first")
    }
    feed.Instance.StartService()

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, err := b.Bot.GetUpdatesChan(u)
    if err != nil {
        panic(err)
    }
    log.Println("Start listening")

    for update := range updates {
        go command.HandleCommand(b.Bot, &update)
    }
}

func main() {
    appConfig := config.LoadConfig("./config.json")
    feedBot := FeedBot{
        Config: &appConfig,
    }
    feedBot.init()
    model.Init(appConfig)
    feed.Instance.Init(appConfig, feedBot.Bot)
    feedBot.run()
}
