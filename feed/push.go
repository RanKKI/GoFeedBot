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

func PushFeedServices(bot *tgbotapi.BotAPI) chan *UserFeeds {
    ch := make(chan *UserFeeds)
    go func(bot *tgbotapi.BotAPI, feeds chan *UserFeeds) {
        for feed := range feeds {
            for _, item := range feed.Items {
                content := cleanContent(item.Description)
                msg := tgbotapi.NewMessage(feed.ChatID, "")
                msg.Text = fmt.Sprintf("<b>%s</b>\n%s\n%s", item.Title, content, item.Link)
                msg.ParseMode = tgbotapi.ModeHTML
                if _, err := bot.Send(msg); err != nil {
                    log.Panicln(err)
                }
            }
        }
    }(bot, ch)
    return ch
}
