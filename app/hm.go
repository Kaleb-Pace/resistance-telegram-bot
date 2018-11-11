package main

import (
	"strings"
)

var hmCommand = BotCommand{
	name:        "Hm",
	description: "Hm",
	matcher: func(update Update) bool {
		return (update.Message.From.UserName != "JacobMason" &&
			(strings.ToLower(update.Message.Text) == "hm" || strings.ToLower(update.Message.Text) == "mm"))
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("© 2018 Jacob Mason", update.Message.Chat.ID)
	},
}
