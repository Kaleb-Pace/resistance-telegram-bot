package main

import (
	"fmt"
	"strconv"
	"strings"
)

var yeetCommand = BotCommand{
	name:        "Yeet",
	description: "Yeetus that Feetus",
	matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "/yeet"
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {

		chadID := strconv.FormatInt(update.Message.Chat.ID, 10)
		userID := strconv.Itoa(update.Message.From.ID)
		yeets, err := usersSpiteLeft(chadID, userID)

		if err != nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("My yeets are all yotted: %s", err.Error()), update.Message.Chat.ID)
			return
		}

		if update.Message.ReplyToMessage == nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("fetusesets: %d", yeets), update.Message.Chat.ID)
			return
		}

		if update.Message.ReplyToMessage.From.UserName == "elicdavis_resistance_bot" {
			respChan <- *NewTextBotResponse("Ya can't yeet da yeetest", update.Message.Chat.ID)
			return
		}

		if yeets <= 0 {
			respChan <- *NewTextBotResponse("This yotter is yeeted out", update.Message.Chat.ID)
		} else {

			err := tradeSpite(chadID, userID, strconv.Itoa(update.Message.ReplyToMessage.From.ID))

			if err == nil {
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s yeeted %s", update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName), update.Message.Chat.ID)
				go bot.deleteMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
				go bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			} else {
				respChan <- *NewTextBotResponse(fmt.Sprintf("My yeets are all yotted: %s", err.Error()), update.Message.Chat.ID)
			}

		}

	},
}
