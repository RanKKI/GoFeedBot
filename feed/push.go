package feed

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
    "strings"
)

var removeTag = regexp.MustCompile("<[^>]*>")
var removeBlank = regexp.MustCompile("\n+\\s*")

func cleanContent(raw string) string {
    output := ""
    for _, breakTag := range []string{"li", "p"} {
        raw = strings.ReplaceAll(raw, fmt.Sprintf("</%s>", breakTag), "\n")
    }
    raw = strings.ReplaceAll(raw, "<li>", "- ")

    raw = removeTag.ReplaceAllString(raw, "")
    raw = removeBlank.ReplaceAllString(raw, "\n")

    for _, val := range strings.Split(raw, "\n") {
        if len(output) > 350 {
            break
        }
        if len(output) > 300 && len(val) > 100 {
            break
        }
        output += "\n" + val
    }

    return output + "......"
}

func PushFeedServices(bot *tgbotapi.BotAPI) chan *UserFeed {
    ch := make(chan *UserFeed)
    go func(bot *tgbotapi.BotAPI, ch chan *UserFeed) {
        for feed := range ch {
            go func(f *UserFeed) {
                content := cleanContent(f.Item.Description)
                msg := tgbotapi.NewMessage(f.ChatID, "")
                msg.Text = fmt.Sprintf("<b>%s</b>\n%s\n%s", f.Item.Title, content, f.Item.Link)
                msg.ParseMode = tgbotapi.ModeHTML
                if _, err := bot.Send(msg); err != nil {
                    log.Panicln(err)
                }
            }(feed)
        }
    }(bot, ch)
    return ch
}
