package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

// TeleBot talks to telegram and manages application state
type TeleBot struct {
	key          string
	lastUpdate   int
	url          string
	chatBuffers  map[string]MessageStack
	chatAliases  map[string]string
	commands     []Command
	botResponses chan BotResponse
	errorReport  Report
	redditUser   RedditAccount
	masterDb     *sql.DB
}

// NewTelegramBot Creates a new telegram bot
func NewTelegramBot(key string, errorReport Report, redditAccount RedditAccount, masterDb *sql.DB, commands []Command) *TeleBot {
	t := TeleBot{
		botResponses: make(chan BotResponse),
		chatAliases:  make(map[string]string),
		chatBuffers:  make(map[string]MessageStack),
		commands:     commands,
		errorReport:  errorReport,
		key:          key,
		lastUpdate:   0,
		redditUser:   redditAccount,
		url:          fmt.Sprintf("https://api.telegram.org/bot%s/", key),
		masterDb:     masterDb,
	}
	return &t
}

func (telebot TeleBot) GetCommands() []Command {
	return telebot.commands
}

func (telebot TeleBot) IsAliasSet(alias string) (string, bool) {
	str, b := telebot.chatAliases[alias]
	return str, b
}

func (telebot TeleBot) SetChatAlias(alias string, chatID int64) {
	telebot.chatAliases[alias] = strconv.FormatInt(chatID, 10)
}

// PushMessageToChatBuffer moves a message to the appropriate chats
func (telebot *TeleBot) PushMessageToChatBuffer(lookup string, message Message) {
	location := strconv.FormatInt(message.Chat.ID, 10)
	if lookup != "" {
		location = lookup
		alias, exists := telebot.chatAliases[location]
		if exists {
			location = alias
		}
	}
	telebot.chatBuffers[location] = telebot.chatBuffers[location].Push(message)
}

// ClearBuffer clears the chat's buffer and returns what has been removed
func (telebot *TeleBot) ClearBuffer(chatID int64) MessageStack {
	lookup := strconv.FormatInt(chatID, 10)
	buffer := telebot.chatBuffers[lookup]
	telebot.chatBuffers[lookup] = make([]Message, 0)
	return buffer
}

// ChatBuffer returns the buffer for a specific chat given the lookup
func (telebot TeleBot) ChatBuffer(lookup string) MessageStack {

	// Try looking it up immediately like they gave us a chat id
	buffer, exists := telebot.chatBuffers[lookup]
	if exists {
		return buffer
	}

	// If we didn't find anything, they might of given us an alias
	alias, exists := telebot.chatAliases[lookup]
	if exists {
		return telebot.chatBuffers[alias]
	}

	return nil
}

// SendMessage Resposible for sending a message to the appropriate group chat
func (telebot TeleBot) sendMessage(message string, chatID int64) {

	postValues := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"text":       message,
		"parse_mode": "HTML",
	}

	jsonValue, err := json.Marshal(postValues)

	if err != nil {
		telebot.errorReport.Log("Error encoding json: " + err.Error())
		return
	}

	req, err := http.NewRequest("POST", telebot.url+"sendMessage", bytes.NewBuffer(jsonValue))

	if err != nil {
		telebot.errorReport.Log("Error creating message: " + err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		telebot.errorReport.Log("Error sending message: " + err.Error())
		return
	}

	// Catch errors
	if err != nil {
		telebot.errorReport.Log("Error sending message: " + err.Error())
		return
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		telebot.errorReport.Log("Error recieving response from tele: " + err.Error())
	}
}

// GetUpdates queries telegram for latest updates
func (telebot *TeleBot) GetUpdates() ([]Update, error) {
	resp, err := http.Get(telebot.url + "getUpdates?offset=" + strconv.Itoa(telebot.lastUpdate))

	// Sometimes Telegram will just randomly send a 502
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	defer resp.Body.Close()

	var updates BatchUpdates
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &updates)
	if err != nil {
		return nil, err
	}

	for _, update := range updates.Result {
		telebot.lastUpdate = update.UpdateID + 1
	}

	return updates.Result, nil
}

func (telebot TeleBot) GetFile(fileID string, byteLimit int) (*http.Response, *GetFileResponseData, error) {
	resp, err := http.Get(fmt.Sprintf("%sgetFile?file_id=%s", telebot.url, fileID))

	log.Println("Begining download")

	if err != nil {
		log.Println("Error: " + err.Error())
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var imageResponse GetFileResponse
	err = json.Unmarshal([]byte(body), &imageResponse)
	if err != nil {
		log.Println("err: " + err.Error())
		return nil, nil, err
	}

	if imageResponse.Ok == false {
		return nil, nil, errors.New("telegram resolved unsucessfully")
	}

	if imageResponse.Result.FileSize > byteLimit {
		return nil, nil, fmt.Errorf("Filesize exceed limit: %d > %d", imageResponse.Result.FileSize, byteLimit)
	}

	log.Println(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", telebot.key, imageResponse.Result.FilePath))

	resp, err = http.Get(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", telebot.key, imageResponse.Result.FilePath))
	return resp, &imageResponse.Result, err
}

// DownloadFile Downloads a file using its file id and returns the filepath on the system
func (telebot TeleBot) DownloadFile(fileID string, byteLimit int) (string, error) {

	resp, responseData, err := telebot.GetFile(fileID, byteLimit)
	if err != nil {
		return "", err
	}

	folder := "media/"
	fileName := responseData.FilePath

	splitResults := strings.Split(responseData.FilePath, "/")

	if len(splitResults) == 2 {
		folder += splitResults[0]
		fileName = splitResults[1]
	}

	os.MkdirAll(folder, os.ModePerm)

	output, err := os.Create(folder + "/" + fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(output, resp.Body)

	if err != nil {
		return "", err
	}

	log.Println("Succesfully downloaded")

	return folder + "/" + fileName, nil
}

func (telebot TeleBot) deleteMessage(chatID int64, messageID int) (bool, error) {

	resp, err := http.Get(fmt.Sprintf("%sdeleteMessage?chat_id=%s&message_id=%d", telebot.url, strconv.FormatInt(chatID, 10), messageID))

	if err != nil {
		log.Println("Error: " + err.Error())
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf(string(body))
	return true, nil

}

// SendPhotoByID send photo by already existing file id
func (telebot TeleBot) SendPhotoByID(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendPhoto?chat_id=%s&photo=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) SendFileByID(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendDocument?chat_id=%s&document=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) SendSticker(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendSticker?chat_id=%s&sticker=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) SendVideoByID(fileID string, chatID int64) error {
	_, err := http.Get(fmt.Sprintf("%ssendvideo?chat_id=%s&video=%s", telebot.url, strconv.FormatInt(chatID, 10), fileID))
	return err
}

func (telebot TeleBot) startForm(chatID int64) (*multipart.Writer, *bytes.Buffer) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	return w, &b
}

func (telebot TeleBot) sendFile(path string, chatID int64) {
	multipartWriter, buffer := telebot.startForm(chatID)

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		telebot.errorReport.Log(err.Error())
	}
	defer file.Close()

	fw, err := multipartWriter.CreateFormFile("document", "movie.mp4")
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		log.Println(err)
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", telebot.url+"sendDocument", buffer)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		telebot.errorReport.Log(err.Error())

		bytes, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Println(err)
			telebot.errorReport.Log(err.Error())
		}
		log.Println(string(bytes))
	}
}

func (telebot TeleBot) SendPhotoByContext(context *gg.Context, chatID int64) {
	multipartWriter, buffer := telebot.startForm(chatID)

	var fw io.Writer

	fw, err := multipartWriter.CreateFormFile("photo", "image.png")

	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	if err := context.EncodePNG(fw); err != nil {
		telebot.errorReport.Log(err.Error())
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", telebot.url+"sendPhoto", buffer)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

}

func (telebot TeleBot) sendImage(path string, chatID int64) {
	multipartWriter, buffer := telebot.startForm(chatID)

	file, err := os.Open(path)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	img, _, err := image.Decode(file)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	fw, err := multipartWriter.CreateFormFile("photo", "image.png")

	if err != nil {
		telebot.errorReport.Log(err.Error())
	}
	if err = png.Encode(fw, img); err != nil {
		telebot.errorReport.Log(err.Error())
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", telebot.url+"sendPhoto", buffer)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		telebot.errorReport.Log(err.Error())
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		telebot.errorReport.Log(err.Error())
	}
	log.Println(string(bytes))
}

func (telebot TeleBot) Start() {
	go func() {
		for response := range telebot.botResponses {
			response.Execute(telebot)
		}
	}()
}

func (telebot *TeleBot) saveMessageToDB(message *Message) {
	if message == nil {
		return
	}

	var replyTo *int
	if message.ReplyToMessage != nil {
		replyTo = &message.ReplyToMessage.MessageID
	}

	var fromID *int
	var fromUserName *string
	if message.From != nil {
		fromID = &message.From.ID
		fromUserName = &message.From.UserName
	}

	var forwardedFromUserID *int
	if message.ForwardFrom != nil {
		forwardedFromUserID = &message.ForwardFrom.ID
	}

	var forwardedFromChatID *int64
	if message.ForwardFromChat != nil {
		forwardedFromChatID = &message.ForwardFromChat.ID
	}

	var photoFileID *string
	if message.Photo != nil {
		photoFileID = &(*message.Photo)[0].FileID
	}

	var videoFileID *string
	if message.Video != nil {
		videoFileID = &message.Video.FileID
	}

	var documentFileID *string
	if message.Document != nil {
		documentFileID = &message.Document.FileID
	}

	var stickerFileID *string
	if message.Sticker != nil {
		stickerFileID = &message.Sticker.FileID
	}

	_, err := telebot.masterDb.Exec(
		"INSERT INTO messages (MessageID, ChatID, Date, FromID, FromUserName, ReplyToMessageID, ForwardedFromUserID, ForwardedFromChatID, PhotoFileID, VideoFileID, DocumentFileID, StickerID, Text) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.MessageID,
		message.Chat.ID,
		message.Date,
		fromID,
		fromUserName,
		replyTo,
		forwardedFromUserID,
		forwardedFromChatID,
		photoFileID,
		videoFileID,
		documentFileID,
		stickerFileID,
		message.Text,
	)

	if err != nil {
		log.Printf("Error saving message: %s\n", err.Error())
	}
}

func (telebot *TeleBot) OnMessage(update Update) {
	go telebot.saveMessageToDB(update.Message)
	for _, command := range telebot.commands {
		if command.Matcher(update) {
			go command.Execute(*telebot, update, telebot.botResponses)
		}
	}
}
