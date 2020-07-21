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

func PushAllFeeds(bot *tgbotapi.BotAPI, userFeeds []*UserFeeds) {
    for _, feeds := range userFeeds {
        for _, item := range feeds.Items {
            content := cleanContent(item.Description)
            msg := tgbotapi.NewMessage(feeds.ChatID, "")
            msg.Text = fmt.Sprintf("<b>%s</b>\n%s\n%s", item.Title, content, item.Link)
            msg.ParseMode = tgbotapi.ModeHTML
            if _, err := bot.Send(msg); err != nil {
                log.Panicln(err)
            }
        }
    }
}
