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
    log.Printf("User %d used command /%s", update.Message.Chat.ID, cmdStr)
    msg := cmd(update)
    if _, err := bot.Send(msg); err != nil {
        log.Panic(err)
    }
}
