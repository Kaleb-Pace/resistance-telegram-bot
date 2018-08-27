package main

import (
	"regexp"
	"strings"
)

var doitCommand = BotCommand{
	Name:        "Do it",
	Description: "You won't. No Balls",
	Matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "do it"
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("You won't", update.Message.Chat.ID)
		respChan <- *NewTextBotResponse("No balls", update.Message.Chat.ID)
	},
}

var wontRegex = regexp.MustCompile("[Yy]ou won'?t")

var youwontCommand = BotCommand{
	Name:        "You wont",
	Description: "No Balls",
	Matcher: func(update Update) bool {
		return wontRegex.FindString(update.Message.Text) == update.Message.Text && update.Message.Text != ""
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("No balls", update.Message.Chat.ID)
	},
}
