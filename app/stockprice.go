package main

import (
	"regexp"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
)

var tickerRegex = regexp.MustCompile("\\$[A-Z]+")

var stockPrice = BotCommand{
	name:        "Stocks",
	description: "Returns the price information for a ticker symbol using iextrading.com",
	matcher: func(update Update) bool {
		return len(tickerRegex.FindAllString(update.Message.Text, -1)) > 0
	},
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		tickers := tickerRegex.FindAllString(update.Message.Text, -1)

		for _, ticker := range tickers {
			url := fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/price", ticker[1:])
			resp, err := http.Get(url)
			if (err == nil && resp.StatusCode == 200) {
				priceText, _ := ioutil.ReadAll(resp.Body)
				var price, _ = strconv.ParseFloat(string(priceText[:]), 64)
				respChan <- *NewTextBotResponse(fmt.Sprintf("%s: %.2f\n", ticker, price), update.Message.Chat.ID)
			}
			resp.Body.Close()
		}
	},
}
