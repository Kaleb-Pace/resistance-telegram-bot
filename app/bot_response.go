package main

import (
	"log"

	"github.com/fogleman/gg"
)

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

/*************************** PictureContextBotResponse ***************************/
type PictureContextBotResponse struct {
	context *gg.Context
	chatID  int64
}

func NewPictureContextBotResponse(context *gg.Context, chatID int64) *PictureContextBotResponse {
	return &PictureContextBotResponse{context, chatID}
}

func (r PictureContextBotResponse) Execute(telebot TeleBot) {
	telebot.SendPhotoByContext(r.context, r.chatID)
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

/*************************** FileReferenceResponse ***************************/
type FlieReferenceBotResponse struct {
	fid    string
	chatID int64
}

func NewFileReferenceBotResponse(fid string, chatID int64) *FlieReferenceBotResponse {
	return &FlieReferenceBotResponse{fid, chatID}
}

func (r FlieReferenceBotResponse) Execute(telebot TeleBot) {
	telebot.SendFileByID(r.fid, r.chatID)
}

/*************************** VideoReferenceResponse ***************************/
type VideoReferenceResponse struct {
	fid    string
	chatID int64
}

func NewVideoReferenceBotResponse(fid string, chatID int64) *VideoReferenceResponse {
	return &VideoReferenceResponse{fid, chatID}
}

func (r VideoReferenceResponse) Execute(telebot TeleBot) {
	err := telebot.SendVideoByID(r.fid, r.chatID)
	if err != nil {
		log.Printf("Error Sending Video: %s", err.Error())
	} else {
		log.Print("Video Sent successfully")
	}
}
