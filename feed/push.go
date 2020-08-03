package feed

import (
    "GoTeleFeed/config"
    "GoTeleFeed/utils"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

type Pusher struct {
    Stream chan *UserFeed
    Bot    *tgbotapi.BotAPI
    Config config.Config
}

func (p *Pusher) StartPushServices() {
    for feed := range p.Stream {
        go func(f *UserFeed) {
            content := utils.CleanHtmlContent(f.Item.Description, p.Config.MaxContentLength)
            msg := tgbotapi.NewMessage(f.ChatID, "")
            msg.Text = fmt.Sprintf("<b>%s</b>\n%s\n%s", f.Item.Title, content, f.Item.Link)
            msg.ParseMode = tgbotapi.ModeHTML
            if _, err := p.Bot.Send(msg); err != nil {
                log.Printf("Error on sending message to %d %s", f.ChatID, err.Error())
            }
        }(feed)
    }
}
