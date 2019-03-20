package main

import (
	"io/ioutil"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

var sunnyCommand = BotCommand{
	name:        "Sunny",
	description: "Generate a Sunny Image",
	matcher:     messageContainsCommandMatcher("sunny"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		wholeCommand := getContentFromCommand(update.Message.Text, "sunny")

		if wholeCommand == "" {
			if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text != "" {
				wholeCommand = update.Message.ReplyToMessage.Text
			} else {
				respChan <- *NewTextBotResponse("Please give me text or respond to a text message", update.Message.Chat.ID)
				return
			}
		}

		if string(wholeCommand[0]) != "\"" {
			wholeCommand = "\"" + wholeCommand
		}

		if string(wholeCommand[len(wholeCommand)-1]) != "\"" {
			wholeCommand = wholeCommand + "\""
		}

		fontByteData, err := ioutil.ReadFile("./TextileRegular.ttf")
		font, err := truetype.Parse(fontByteData)
		if err != nil {
			bot.errorReport.Log(err.Error())
			respChan <- *NewTextBotResponse("Error Loading font: "+err.Error(), update.Message.Chat.ID)
			return
		}

		imageWidth := 400
		imageHeight := 300

		context := gg.NewContext(imageWidth, imageHeight)
		context.SetRGB(0, 0, 0)
		context.DrawRectangle(0, 0, float64(imageWidth), float64(imageHeight))
		context.Fill()
		context.SetRGB(1, 1, 1)

		fontSize := 60.0

		var lines []string = nil

		for lines == nil || float64(len(lines))*fontSize > float64(imageHeight)*.2 {
			fontSize *= .9
			face := truetype.NewFace(font, &truetype.Options{
				Size: fontSize,
			})

			context.SetFontFace(face)

			lines = context.WordWrap(wholeCommand, float64(imageWidth)*0.8)
		}

		totalHeight := 0.0
		for _, l := range lines {
			_, h := context.MeasureString(l)
			totalHeight += h
		}

		heightAddedSoFar := 0.0
		for _, l := range lines {
			w, h := context.MeasureString(l)
			context.DrawString(l, (float64(imageWidth)-w)/2.0, ((float64(imageHeight)-totalHeight)/2.0)+heightAddedSoFar+(h/2))
			heightAddedSoFar += h
		}

		respChan <- *NewPictureContextBotResponse(context, update.Message.Chat.ID)
	},
}
