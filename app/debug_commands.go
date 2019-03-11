package main

import (
	"fmt"
	"strings"
)

var userIDCommand = BotCommand{
	name:        "userId",
	description: "see the user's id",
	matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "/uid"
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		if update.Message.ReplyToMessage == nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("User Id: %d", update.Message.From.ID), update.Message.Chat.ID)
		} else {
			respChan <- *NewTextBotResponse(fmt.Sprintf("User Id: %d", update.Message.ReplyToMessage.From.ID), update.Message.Chat.ID)
		}
	},
}

var messageIDCommand = BotCommand{
	name:        "messageId",
	description: "see the user's id",
	matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "/mid"
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		if update.Message.ReplyToMessage != nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("Message Id: %d", update.Message.ReplyToMessage.MessageID), update.Message.Chat.ID)
		}
	},
}
