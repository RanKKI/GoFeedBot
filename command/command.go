package command

import (
    "GoTeleFeed/database"
    "GoTeleFeed/feed"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

type command func(update *tgbotapi.Update) tgbotapi.Chattable

var commands = map[string]command{
    "start": cmdStart,
    "add":   cmdAdd,
    "list":  cmdList,
    "help":  cmdHelp,
}

func cmdHelp(update *tgbotapi.Update) tgbotapi.Chattable {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
    msg.ReplyToMessageID = update.Message.MessageID

    msg.Text = "/add <url> # Subscribe a new feed url"
    msg.Text += "\n/list # Show all Subscribed URLs"
    msg.Text += "\n/help # Show this message"

    return msg
}

func cmdStart(update *tgbotapi.Update) tgbotapi.Chattable {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello~ My friend")
    msg.ReplyToMessageID = update.Message.MessageID
    return msg
}

func cmdAdd(update *tgbotapi.Update) tgbotapi.Chattable {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
    msg.ReplyToMessageID = update.Message.MessageID

    url := update.Message.CommandArguments()

    f, err := feed.TestFeed(url)

    if err != nil {
        msg.Text = "Failed, " + err.Error()
        return msg
    }

    log.Printf("User %d, Subsribe %s", msg.ChatID, url)

    // 如果成功, 保存进数据库

    if err := database.AddFeed(msg.ChatID, f, url); err != nil {
        msg.Text = err.Error()
        return msg
    }
    msg.Text = fmt.Sprintf("Success.\nTitle: %s", f.Title)
    if f.Author != nil {
        msg.Text += "\nby: " + f.Author.Name
    }
    return msg
}

func cmdList(update *tgbotapi.Update) tgbotapi.Chattable {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
    msg.ReplyToMessageID = update.Message.MessageID

    for _, f := range database.QueryFeeds(msg.ChatID) {
        msg.Text += fmt.Sprintf("\n%d - %s", f.ID, f.Title)
    }

    return msg
}
