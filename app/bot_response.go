package main

// BotResponse set by bot when a command is satisfied
type BotResponse interface {
	Execute(telebot TeleBot)
}

/*************************** TextBotResponse ***************************/
type TextBotResponse struct {
	text   string
	chatID int64
}

func NewTextBotResponse(msg string, chatID int64) *TextBotResponse {
	return &TextBotResponse{msg, chatID}
}

func (r TextBotResponse) Execute(telebot TeleBot) {
	telebot.sendMessage(r.text, r.chatID)
}

/*************************** PictureReferenceBotResponse ***************************/
type PictureReferenceBotResponse struct {
	pid    string
	chatID int64
}

func NewPictureReferenceBotResponse(pid string, chatID int64) *PictureReferenceBotResponse {
	return &PictureReferenceBotResponse{pid, chatID}
}

func (r PictureReferenceBotResponse) Execute(telebot TeleBot) {
	telebot.SendPhotoByID(r.pid, r.chatID)
}

/*************************** PictureUploadBotResponse ***************************/
type PictureUploadBotResponse struct {
	filepath string
	chatID   int64
}

func NewPictureUploadBotResponse(filepath string, chatID int64) *PictureUploadBotResponse {
	return &PictureUploadBotResponse{filepath, chatID}
}

func (r PictureUploadBotResponse) Execute(telebot TeleBot) {
	telebot.sendImage(r.filepath, r.chatID)
}

/*************************** StickerBotResponse ***************************/
type StickerBotResponse struct {
	sid    string
	chatID int64
}

func NewStickerBotResponse(sid string, chatID int64) *StickerBotResponse {
	return &StickerBotResponse{sid, chatID}
}

func (r StickerBotResponse) Execute(telebot TeleBot) {
	telebot.SendSticker(r.sid, r.chatID)
}

/*************************** FileUploadBotResponse ***************************/
type FileUploadBotResponse struct {
	filePath string
	chatID   int64
}

func NewFileUploadBotResponse(filePath string, chatID int64) *FileUploadBotResponse {
	return &FileUploadBotResponse{filePath, chatID}
}

func (r FileUploadBotResponse) Execute(telebot TeleBot) {
	telebot.sendFile(r.filePath, r.chatID)
}
