package main

import (
	"fmt"
	"time"
)

var epochCommand = BotCommand{
	name:        "epoch",
	description: "get current epoch",
	matcher:     messageContainsCommandMatcher("epoch"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse(fmt.Sprintf("%d", time.Now().Unix()), update.Message.Chat.ID)
	},
}
