package main

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var threatCommand = BotCommand{
	name:        "Dicc",
	description: "You know what's coming",
	matcher:     messageContainsCommandMatcher("dicc"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		wholeCommand := getContentFromCommand(update.Message.Text, "dicc")

		if wholeCommand == "" {
			if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text != "" {
				wholeCommand = update.Message.ReplyToMessage.Text
			} else {
				return
			}
		}

		font, err := truetype.Parse(goregular.TTF)
		if err != nil {
			bot.errorReport.Log(err.Error())
		}

		im, err := gg.LoadPNG("picturecommands/kable-in.png")
		if err != nil {
			bot.errorReport.Log("unable to load image: " + err.Error())
			return
		}
		dc := gg.NewContextForImage(im)
		dc.SetRGB(0, 0, 0)

		dc.SetRGB(1, 1, 1)

		face := truetype.NewFace(font, &truetype.Options{
			Size: 50,
		})

		widthOffset := float64(len(wholeCommand)/2) * 20.0
		dc.SetFontFace(face)
		dc.RotateAbout(-math.Pi*.30, 130, 600)
		dc.DrawString(wholeCommand, 150-widthOffset, 650)

		respChan <- *NewPictureContextBotResponse(dc, update.Message.Chat.ID)
	},
}
