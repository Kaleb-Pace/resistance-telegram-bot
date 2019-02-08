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
			return
		}

		font, err := truetype.Parse(goregular.TTF)
		if err != nil {
			bot.errorReport.Log(err.Error())
		}
		face := truetype.NewFace(font, &truetype.Options{
			Size: 70,
		})

		im, err := gg.LoadPNG("picturecommands/sarahsays.png")
		if err != nil {
			bot.errorReport.Log("unable to load image: " + err.Error())
			return
		}
		dc := gg.NewContextForImage(im)
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(face)
		dc.DrawStringAnchored(update.Message.From.UserName, 400, 100, 0.0, 0.0)
		dc.DrawStringAnchored(wholeCommand, 100, 100, 0.0, 0.0)
		dc.SavePNG("sarahout.png")

		respChan <- *NewFileBotResponse("sarahout.png", update.Message.Chat.ID)
	},
}
