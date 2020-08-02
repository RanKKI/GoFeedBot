package command

import (
    "GoTeleFeed/feed"
    "GoTeleFeed/model"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strconv"
)

type command func(update *tgbotapi.Update) tgbotapi.MessageConfig

var commands = map[string]command{
    "start":  cmdStart,
    "add":    cmdAdd,
    "list":   cmdList,
    "help":   cmdHelp,
    "remove": cmdRemove,
}

func cmdHelp(update *tgbotapi.Update) tgbotapi.MessageConfig {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

    msg.Text = "/add <url> # Subscribe a new feed url\n" +
        "/remove <id> # Remove a exists subscribe\n" +
        "/list # Show all subscribes\n" +
        "/help # Show this message"

    return msg
}

func cmdStart(update *tgbotapi.Update) tgbotapi.MessageConfig {
    return tgbotapi.NewMessage(update.Message.Chat.ID, "Hello~ My friend")
}

func cmdAdd(update *tgbotapi.Update) tgbotapi.MessageConfig {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

    url := update.Message.CommandArguments()

    if len(url) == 0 {
        msg.Text = "Usage: /add <feed_url>"
        return msg
    }

    title := ""
    name := ""

    // Since there is a same feed url in the db
    // No need to check whether is a valid link
    if fe := model.QueryFeed(model.Feed{URL: url}); fe != nil {
        title = fe.Title
        name = fe.AuthorName

    } else {
        f, err := feed.Instance.TestURL(url)

        if err != nil {
            msg.Text = "Failed, " + err.Error()
            return msg
        }

        title = f.Title
        if f.Author != nil {
            name = f.Author.Name
        }

    }

    if err := model.AddSub(msg.ChatID, url, title, name); err != nil {
        msg.Text = err.Error()
        return msg
    }

    log.Printf("User %d subsribed %s", msg.ChatID, url)

    msg.Text = fmt.Sprintf("Success.\nTitle: %s", title)
    if len(name) != 0 {
        msg.Text += "\nby: " + name
    }
    return msg
}

func cmdList(update *tgbotapi.Update) tgbotapi.MessageConfig {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

    feeds := model.QueryFeeds(model.Feed{ChatID: msg.ChatID})
    if len(feeds) == 0 {
        msg.Text = "You don't have any subscribes at the moment"
        return msg
    }
    for _, f := range feeds {
        msg.Text += fmt.Sprintf("\n%d - %s", f.ID, f.Title)
    }

    return msg
}

func cmdRemove(update *tgbotapi.Update) tgbotapi.MessageConfig {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
    arg := update.Message.CommandArguments()

    if len(arg) == 0 {
        msg.Text = "Usage: /remove <feed_id>"
        return msg
    }

    subID, err := strconv.ParseInt(arg, 10, 64)
    if err != nil {
        msg.Text = fmt.Sprintf("Invalid argument %s", arg)
        return msg
    }

    err = model.DeleteSub(update.Message.Chat.ID, subID)

    if err != nil {
        msg.Text = err.Error()
        return msg
    }

    msg.Text = "Succeed"

    return msg
}
