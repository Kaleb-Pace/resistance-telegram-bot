package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type QueryCommand struct {
	command *sql.DB
}

func NewQueryCommand(command *sql.DB) *QueryCommand {
	return &QueryCommand{command}
}

// Name returns name of the command
func (b QueryCommand) Name() string {
	return "Query"
}

// Description returns the description of the command
func (b QueryCommand) Description() string {
	return "Explore the fruits of datamining. Start the query with /sql"
}

// Matcher determines whether or not the execute will commense
func (b QueryCommand) Matcher(update Update) bool {
	return len(strings.SplitAfter(update.Message.Text, "/sql")) > 1
}

// Execute runs the command logic
func (b QueryCommand) Execute(telebot TeleBot, update Update, respChan chan BotResponse) {
	query := getContentFromCommand(update.Message.Text, "sql")

	sqlQuery := fmt.Sprintf("SELECT %s", query)

	log.Println(sqlQuery)

	rows, err := telebot.masterDb.Query(sqlQuery)

	if err != nil {
		respChan <- *NewTextBotResponse(fmt.Sprintf("Error making query: %s", err.Error()), update.Message.Chat.ID)
		return
	}

	finalMessage := ""

	defer rows.Close()
	cols, _ := rows.Columns()

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("Error scanning row: %s", err.Error()), update.Message.Chat.ID)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		for i, c := range result {
			finalMessage += fmt.Sprintf("%v ", c)
			if i != len(result)-1 {
				finalMessage += "- "
			}
		}

		finalMessage += "\n"
	}

	respChan <- *NewTextBotResponse(finalMessage, update.Message.Chat.ID)

}
