package main

import (
	"regexp"
)

var resistanceRuleOneRegex = regexp.MustCompile("(R|r)ule #?1$")
var resistanceRuleTwoRegex = regexp.MustCompile("(R|r)ule #?2$")
var resistanceRuleThreeRegex = regexp.MustCompile("(R|r)ule #?3$")

var resistanceRuleOneCommand = BotCommand{
	name:        "Rule 1",
	description: "We don't appreciate creativity or talent here",
	matcher: func(update Update) bool {
		return len(resistanceRuleOneRegex.FindAllString(update.Message.Text, -1)) > 0
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("We don't appreciate creativity or talent here", update.Message.Chat.ID)
	},
}

var resistanceRuleTwoCommand = BotCommand{
	name:        "Rule 2",
	description: "NSFW After 4:59pmCST and on Weekends",
	matcher: func(update Update) bool {
		return len(resistanceRuleTwoRegex.FindAllString(update.Message.Text, -1)) > 0
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("NSFW After 4:59pmCST and on Weekends", update.Message.Chat.ID)
	},
}

var resistanceRuleThreeCommand = BotCommand{
	name:        "Rule 3",
	description: "Merge Request or STFU",
	matcher: func(update Update) bool {
		return len(resistanceRuleThreeRegex.FindAllString(update.Message.Text, -1)) > 0
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		respChan <- *NewTextBotResponse("Merge Request or STFU", update.Message.Chat.ID)
	},
}
