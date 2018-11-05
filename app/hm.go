package main

import (
	"strings"
)

var hmCommand = BotCommand{
	Name:        "Hm",
	Description: "Hm",
	Matcher: func(update Update) bool {
		return (update.Message.From.UserName != "JacobMason" &&
			(strings.ToLower(update.Message.Text) == "hm" || strings.ToLower(update.Message.Text) == "mm"))
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("© 2018 Jacob Mason", update.Message.Chat.ID)
	},
}
