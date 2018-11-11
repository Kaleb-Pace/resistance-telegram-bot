package main

// Command is a type of bot action that can be executed on a message
type Command interface {
	Matcher(msg Update) bool
	Execute(x TeleBot, y Update, w chan BotResponse)
	Name() string
	Description() string
}
