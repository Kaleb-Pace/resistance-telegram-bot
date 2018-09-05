package main

import (
	"fmt"
	"log"
)

var yeetCommand = BotCommand{
	Name:        "Yeet",
	Description: "Yeetus that Feetus",
	Matcher:     messageContainsCommandMatcher("yeet"),
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		log.Print(update.Message.ReplyToMessage.From.UserName)
		if update.Message.ReplyToMessage.From.UserName == "elicdavis_resistance_bot" {
			respChan <- *NewTextBotResponse("Ya can't yeet da yeetest", update.Message.Chat.ID)
		} else {
			respChan <- *NewTextBotResponse(fmt.Sprintf("%s yeeted %s", update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
			go bot.deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
			go bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		}
	},
}
