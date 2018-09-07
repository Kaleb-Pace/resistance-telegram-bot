package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func usersYeetsLeft(chatId string, poster string) (int, error) {
	f, err := os.OpenFile("yeets/"+chatId, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		if poster == entry[0] {
			return strconv.Atoi(entry[1])
		}
	}

	return 3, nil
}

type yeetEntry struct {
	user  string
	yeets int
}

func tradeYeets(chatID string, from string, to string) error {
	f, err := os.OpenFile("yeets/"+chatID, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	yeeters := make([]yeetEntry, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		yeets, err := strconv.Atoi(entry[1])
		if err != nil {
			return err
		}
		yeeters = append(yeeters, yeetEntry{
			user:  entry[0],
			yeets: yeets,
		})
		y := yeeters[len(yeeters)-1]
		log.Printf("[%s]:%d\n", y.user, y.yeets)
	}

	yeeterFound := false
	yeetedFound := false
	for i, yeeter := range yeeters {
		if yeeter.user == from {
			yeeters[i].yeets--
			yeeterFound = true
			log.Printf("found yeeter")
		} else if yeeter.user == to {
			yeeters[i].yeets++
			yeetedFound = true
			log.Printf("found yeeted")
		}
	}

	if yeeterFound == false {
		yeeters = append(yeeters, yeetEntry{
			user:  from,
			yeets: 2,
		})
	}

	if yeetedFound == false {
		yeeters = append(yeeters, yeetEntry{
			user:  to,
			yeets: 4,
		})
	}

	f.Truncate(0)
	f.Seek(0, 0)

	for _, yeeter := range yeeters {
		if _, err = f.WriteString(fmt.Sprintf("%s %d\n", yeeter.user, yeeter.yeets)); err != nil {
			return err
		}
	}

	return nil
}

var yeetCommand = BotCommand{
	Name:        "Yeet",
	Description: "Yeetus that Feetus",
	Matcher: func(update Update) bool {
		return strings.ToLower(update.Message.Text) == "/yeet"
	},
	Execute: func(bot TeleBot, update Update, respChan chan BotResponse) {

		chadID := strconv.FormatInt(update.Message.Chat.ID, 10)
		userID := strconv.Itoa(update.Message.From.ID)
		yeets, err := usersYeetsLeft(chadID, userID)

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

			err := tradeYeets(chadID, userID, strconv.Itoa(update.Message.ReplyToMessage.From.ID))

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
