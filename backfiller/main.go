package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Export struct {
	Chats chats `json:"chats" binding:"required"`
}

type chats struct {
	List []chat `json:"list"` // optional
}

type chat struct {
	Messagess []message `json:"messages"` // optional
	Name      string    `json:"name"`     // optional
	ID        int64     `json:"id"`       // optional
}

type message struct {
	ID               int    `json:"id"`                  // optional
	MessageType      string `json:"type"`                // optional
	FromID           int    `json:"from_id"`             // optional
	Text             string `json:"text"`                // optional
	ReplyToMessageID int    `json:"reply_to_message_id"` // optional
	Date             string `json:"date"`                // optional
	Photo            string `json:"photo"`               // optional
	ForwardedFrom    string `json:"forwarded_from"`      // optional
	MediaType        string `json:"media_type"`          // optional
	File             string `json:"file"`                // optional
}

var idAliases = map[int]int{
	602338572: 666289914,
}

var idToUsernameMapping = map[int]string{
	129485665: "wesstr",
	398589153: "therealrbp",
	275448822: "jdaltonchilders",
	167815129: "alxxlc",
	240194294: "Austinp96",
	106468411: "B02s2",
	481416215: "BabysMommy",
	106672881: "blondinblue",
	709033588: "fishy_ass",
	365451200: "gandalftheblake",
	370914753: "JacobMason",
	666289914: "Jadec137A",
	346762390: "KevinBall",
	134753805: "KlarkKable",
	405236436: "MarissaHargrave",
	432029462: "nathantfrank",
	302666714: "SlurpChan",
	325303234: "SnarkyPuppy",
	352574940: "TheAlbinoRhino",
}

func (m message) Epoch() int64 {
	t, err := time.Parse("2006-01-02T15:04:05", m.Date)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}

func (m message) ActualReplyTo() string {
	if m.ReplyToMessageID <= 0 {
		return "\\N"
	}
	return strconv.Itoa(m.ReplyToMessageID)
}

func (m message) PhotoID() string {
	if m.Photo != "" {
		return "something"
	}
	return "\\N"
}

func (m message) VideoID() string {
	if m.MediaType == "video_file" || m.MediaType == "animation" {
		return "something"
	}
	return "\\N"
}

func (m message) DocumentID() string {
	if m.File != "" && m.PhotoID() == "\\N" && m.VideoID() == "\\N" && m.StickerID() == "\\N" {
		return "something"
	}
	return "\\N"
}

func (m message) StickerID() string {
	if m.MediaType == "sticker" {
		return "something"
	}
	return "\\N"
}

func (m message) Username() string {
	i, ok := idToUsernameMapping[m.RealID()]
	if !ok {
		panic(fmt.Sprintf("Couldn't find user with id: %d; message id: %d", m.FromID, m.ID))
	}
	return i
}

func (m message) RealID() int {
	i, ok := idAliases[m.FromID]
	if !ok {
		return m.FromID
	}
	return i
}

func (m message) FormattedText() string {
	if m.Text == "" {
		return "\\N"
	}

	cleaned := strings.Replace(m.Text, "\\", "", -1)
	cleaned = strings.Replace(cleaned, ",", "", -1)
	cleaned = strings.Replace(cleaned, "\n", "", -1)

	if cleaned == "" {
		cleaned = "Junk Message Eli Removed When Parsing Exported Data"
	}

	return cleaned
}

func (m message) FormattedForwardedFrom() string {
	if m.ForwardedFrom == "" {
		return "\\N"
	}
	return strings.Replace(m.ForwardedFrom, ",", "", -1)
}

func (m message) CSV(chatID int64) []string {
	return []string{
		strconv.Itoa(m.ID),
		strconv.FormatInt(chatID, 10),
		strconv.FormatInt(m.Epoch(), 10),
		strconv.Itoa(m.RealID()),
		m.Username(),
		m.ActualReplyTo(),
		m.FormattedForwardedFrom(),
		m.PhotoID(),
		m.VideoID(),
		m.DocumentID(),
		m.StickerID(),
		m.FormattedText(),
	}
}

func (m message) Check() error {
	if m.FromID == 0 {
		return errors.New("Unable to get FromID")
	}
	return nil
}

func main() {

	dat, err := ioutil.ReadFile("result-cleaned.json")
	if err != nil {
		panic(err)
	}

	upToMessageID := 200000 //145781
	var chatToPullFrom int64 = 9771743476
	var chatIDAlias int64 = -1001181808884

	usersToIgnore := []int{394077543, 117099167, 469361938}

	var results Export
	json.Unmarshal(dat, &results)

	outFile, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write([]string{
		"messageID",
		"chatID",
		"date",
		"userID",
		"username",
		"ReplyTo",
		"ForwardedFrom",
		"photo",
		"video",
		"document",
		"sticker",
		"text",
	})

	for _, chat := range results.Chats.List {
		// log.Printf("[%d] %s: %d messages", chat.ID, chat.Name, len(chat.Messagess))

		if chat.ID == chatToPullFrom {
			for _, message := range chat.Messagess {

				if message.MessageType != "message" {
					continue
				}

				if message.Check() != nil {
					panic(message.Check())
				}

				if message.ID >= upToMessageID {
					break
				}

				shouldIgnore := false
				for _, id := range usersToIgnore {
					if message.FromID == id {
						shouldIgnore = true
						break
					}
				}

				if shouldIgnore {
					continue
				}

				writer.Write(message.CSV(chatIDAlias))
			}

		}

	}

}
