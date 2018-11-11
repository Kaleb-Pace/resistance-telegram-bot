package main

import (
	"regexp"
	"strings"
)

var doitCommand = BotCommand{
	name:        "Do it",
	description: "You won't. No Balls",
	matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "do it"
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("You won't", update.Message.Chat.ID)
		respChan <- *NewTextBotResponse("No balls", update.Message.Chat.ID)
	},
}

var wontRegex = regexp.MustCompile("[Yy]ou won'?t")

var youwontCommand = BotCommand{
	name:        "You wont",
	description: "No Balls",
	matcher: func(update Update) bool {
		return wontRegex.FindString(update.Message.Text) == update.Message.Text && update.Message.Text != ""
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("No balls", update.Message.Chat.ID)
	},
}
