package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

var wastedCommand = BotCommand{
	name:        "wasted",
	description: "Generated a 'wasted' version of a gif",
	matcher:     messageContainsCommandMatcher("wasted"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {

		if update.Message.ReplyToMessage.Document != nil {
			path, err := bot.GetFile(update.Message.ReplyToMessage.Document.FileID)
			if err != nil {
				respChan <- *NewTextBotResponse("Error downloading file: "+err.Error(), update.Message.Chat.ID)
				return
			}

			fi, err := os.Stat(path)
			if err != nil {
				respChan <- *NewTextBotResponse("Error examining file: "+err.Error(), update.Message.Chat.ID)
				return
			}

			size := fi.Size()
			if size > 2097152 {
				respChan <- *NewTextBotResponse("That file is a little to big for my taste (>2MB). Consider donating vbucks to server maintainer ;3", update.Message.Chat.ID)
				return
			}

			if len(path) <= 4 || path[len(path)-4:] != ".mp4" {
				respChan <- *NewTextBotResponse(fmt.Sprintf("Unsupported filetype: %s", path), update.Message.Chat.ID)
			}

			frameRate, err := VideoFramerate(path)
			if err != nil {
				respChan <- *NewTextBotResponse(fmt.Sprintf("Error retrieving framerate: %s", err.Error()), update.Message.Chat.ID)
				return
			}

			os.MkdirAll("wastedTemp", os.ModePerm)
			RemoveContents("wastedTemp")
			err = UnstitchImages(path, "wastedTemp")

			if err != nil {
				respChan <- *NewTextBotResponse("Error unstitching file: "+err.Error(), update.Message.Chat.ID)
				return
			}

			font, err := truetype.Parse(goregular.TTF)
			if err != nil {
				bot.errorReport.Log(err.Error())
			}

			frames, _ := ioutil.ReadDir("wastedTemp")
			for i := 0; i < len(frames); i++ {
				filePath := fmt.Sprintf("wastedTemp/%04d.png", i+1)
				im, err := gg.LoadPNG(filePath)
				if err != nil {
					bot.errorReport.Log("unable to load image: " + err.Error())
					return
				}

				dc := gg.NewContextForImage(im)
				Greyscale(dc, math.Min(float64(i*2)/float64(len(frames)), 1))
				if float64(i)/float64(len(frames)) > 0.5 {
					size := math.Min(float64(dc.Width()), float64(dc.Height()))
					face := truetype.NewFace(font, &truetype.Options{
						Size: size * .15,
					})
					dc.SetFontFace(face)
					dc.SetRGB(1, 0, 0)
					dc.DrawStringAnchored("WASTED", float64(dc.Width()/2)-(size*.25), float64(dc.Height()/2), 0.0, 0.0)
				}
				dc.SavePNG(filePath)
			}

			err = StichPicturesTogether("wastedTemp/%04d.png", "wastedout.mp4", frameRate)
			if err != nil {
				respChan <- *NewTextBotResponse("Error stitching file: "+err.Error(), update.Message.Chat.ID)
				return
			}

			respChan <- *NewFileUploadBotResponse("wastedout.mp4", update.Message.Chat.ID)

		} else {
			respChan <- *NewTextBotResponse("Please respond to a gif ", update.Message.Chat.ID)
		}

	},
}
