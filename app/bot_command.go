package main

// BotCommand is a basic instance of a command
type BotCommand struct {
	matcher     func(msg Update) bool
	execute     func(x TeleBot, y Update, w chan BotResponse)
	name        string
	description string
}

// Name returns name of the command
func (b BotCommand) Name() string {
	return b.name
}

// Description returns the description of the command
func (b BotCommand) Description() string {
	return b.description
}

// Matcher determines whether or not the execute will commense
func (b BotCommand) Matcher(msg Update) bool {
	return b.matcher(msg)
}

// Execute runs the command logic
func (b BotCommand) Execute(x TeleBot, y Update, w chan BotResponse) {
	b.execute(x, y, w)
}
