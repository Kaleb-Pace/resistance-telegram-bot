package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// DictionarySearchResponse
type DictionarySearchResponse struct {
	RoleInSentence   string   `json:"fl"`
	ShortDefinitions []string `json:"shortdef"`
}

var defineCommand = BotCommand{
	name:        "Define",
	description: "Grabs the definition of a word from Merriam Webster",
	matcher:     messageContainsCommandMatcher("define"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		term := getContentFromCommand(update.Message.Text, "define")

		if term == "" {
			return
		}

		searchURL := fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s", term, os.Getenv("DICT_KEY"))
		resp, err := http.Get(searchURL)

		if err != nil {
			bot.errorReport.Log("Error Searching Dictionary: " + err.Error())
			respChan <- *NewTextBotResponse("Error Searching Dictionary", update.Message.Chat.ID)
			return
		}
		log.Print()

		defer resp.Body.Close()

		var dictionaryResponses []DictionarySearchResponse
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &dictionaryResponses)
		if err != nil {
			bot.errorReport.Log("Error Parsing Dictionary Response: " + err.Error())
			respChan <- *NewTextBotResponse("Error Reading Response From Dictionary", update.Message.Chat.ID)
			return
		}

		var returnMsg bytes.Buffer

		for _, response := range dictionaryResponses {
			returnMsg.WriteString(fmt.Sprintf("\n<b>%s</b>\n", response.RoleInSentence))
			for i, shortDefinition := range response.ShortDefinitions {
				returnMsg.WriteString(fmt.Sprintf("%d. %s\n", i+1, shortDefinition))
			}
		}
		respChan <- *NewTextBotResponse(returnMsg.String(), update.Message.Chat.ID)
	},
}
