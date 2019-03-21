package main

import (
	"database/sql"
	"fmt"
	"log"
)

type QueryCommand struct {
	BotCommand
	database *sql.DB
}

func NewTallyCommand(database *sql.DB) *QueryCommand {
	return &QueryCommand{
		BotCommand: BotCommand{
			name:        "Tally",
			description: "Figure out how often something's been said",
			matcher:     messageContainsCommandMatcher("tally"),
			execute: func(telebot TeleBot, update Update, respChan chan BotResponse) {
				sqlRegex := getContentFromCommand(update.Message.Text, "tally")

				if sqlRegex == "" {
					respChan <- *NewTextBotResponse("Please provide SQL regex", update.Message.Chat.ID)
					return
				}

				result, err := telebot.masterDb.Query(
					"SELECT COUNT(MessageID), FromUserName FROM messages WHERE Text LIKE ? AND ChatID = ? GROUP BY FromUserName",
					sqlRegex,
					update.Message.Chat.ID,
				)

				if err != nil {
					respChan <- *NewTextBotResponse(fmt.Sprintf(err.Error()), update.Message.Chat.ID)
					return
				}

				var amount int
				var username string
				var returnMessage string
				for result.Next() {
					err = result.Scan(&amount, &username)
					if err != nil {
						respChan <- *NewTextBotResponse(fmt.Sprintf(err.Error()), update.Message.Chat.ID)
						return
					}
					returnMessage += fmt.Sprintf("%d - %s\n", amount, username)
				}

				if returnMessage == "" {
					respChan <- *NewTextBotResponse("Doesn't look like any messages text match that sql regex", update.Message.Chat.ID)
				} else {
					respChan <- *NewTextBotResponse(returnMessage, update.Message.Chat.ID)
				}

			},
		},
	}
}

func NewSelectCommand(database *sql.DB) *QueryCommand {
	return &QueryCommand{
		BotCommand: BotCommand{
			name:        "Query",
			description: "Explore the fruits of datamining. Start the query with /sql",
			matcher:     messageContainsCommandMatcher("sql"),
			execute: func(telebot TeleBot, update Update, respChan chan BotResponse) {
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

				dest := make([]interface{}, len(cols))
				for i := range rawResult {
					dest[i] = &rawResult[i]
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

			},
		},
		database: database,
	}
}
