package main

var swallowCommand = BotCommand{
	name:        "Swallow",
	description: "Delete all messages in the buffer",
	matcher:     messageContainsCommandMatcher("swallow"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		buffer := bot.ClearBuffer(update.Message.Chat.ID)
		returnMessage := "No Messages :/"
		if buffer.Size() > 0 {
			returnMessage = "Thanks Daddy :) Messages all gone"
		}
		respChan <- *NewTextBotResponse(returnMessage, update.Message.Chat.ID)
	},
}
