package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type MessageFilter struct {
	ChatID string `json:"chatId"`
}

type FindMessageRequest struct {
	Filter *MessageFilter    `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindMessageResponse struct {
	Messages []*Message `json:"messages"`
	Total    int64      `json:"total"`
}
type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
type PayloadAttachment struct {
	ID          string       `json:"id"`
	Thumbnail   string       `json:"thumbnail"`
	URL         string       `json:"url"`
	Description string       `json:"description"`
	Coordinates *Coordinates `json:"coordinates"`
	Size        string       `json:"size"`
	Name        string       `json:"name"`
	Checksum    string       `json:"checksum"`
	Type        string       `json:"type"`
}
type Attachment struct {
	Payload *PayloadAttachment `json:"payload"`
	Type    string             `json:"type"`
}
type Message struct {
	ID                string        `json:"id"`
	ChatID            string        `json:"chatId"`
	ChatTitle         string        `json:"chatTitle"`
	AuthorID          string        `json:"authorId"`
	Message           string        `json:"message"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
	Attachments       []*Attachment `json:"attachments"`
	ExternalMessageID string        `json:"externalMessageId"`
}
