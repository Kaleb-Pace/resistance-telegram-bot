package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var sarahSaysCommand = BotCommand{
	name:        "Sarah Says",
	description: "Genny is the president of ACMW",
	matcher:     messageContainsCommandMatcher("sarahsays"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {
		wholeCommand := getContentFromCommand(update.Message.Text, "sarahsays")

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

		im, err := gg.LoadPNG("picturecommands/sarahsays.png")
		if err != nil {
			bot.errorReport.Log("unable to load image: " + err.Error())
			return
		}
		dc := gg.NewContextForImage(im)
		dc.SetRGB(0, 0, 0)

		fontSize := 60.0

		var lines []string = nil

		for lines == nil || float64(len(lines))*fontSize > 150.0 {
			fontSize *= .9
			face := truetype.NewFace(font, &truetype.Options{
				Size: fontSize,
			})

			dc.SetFontFace(face)

			lines = dc.WordWrap(wholeCommand, 300)
		}

		for i, l := range lines {
			dc.DrawString(l, 60, 410+(float64(i)*fontSize))
		}

		respChan <- *NewPictureContextBotResponse(dc, update.Message.Chat.ID)
	},
}
