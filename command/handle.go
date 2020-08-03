package command

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
    if update.Message == nil || !update.Message.IsCommand() {
        return
    }

    cmdStr := update.Message.Command()
    cmd, ok := commands[cmdStr]
    if !ok {
        return
    }
    log.Printf("User %d used command %s", update.Message.Chat.ID, update.Message.Text)
    msg := cmd(update)
    msg.ReplyToMessageID = update.Message.MessageID
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Error on repling message to %d %s", msg.ChatID, err.Error())
    }
}
