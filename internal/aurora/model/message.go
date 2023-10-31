package model

import (
	"database/sql"
	"time"
)

const (
	MessageFieldID                = "id"
	MessageFieldChatID            = "chat_id"
	MessageFieldAuthorID          = "author_id"
	MessageFieldMessage           = "message"
	MessageFieldCreatedAt         = "created_at"
	MessageFieldUpdatedAt         = "updated_at"
	MessageFieldDeletedAt         = "deleted_at"
	MessageFieldAttachments       = "attachments"
	MessageFieldExternalMessageID = "external_message_id"
)

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
	ID                string         `db:"id"`
	ChatID            string         `db:"chat_id"`
	AuthorID          sql.NullString `db:"author_id"`
	Message           string         `db:"message"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	Attachments       []*Attachment  `db:"attachments"`
	ExternalMessageID sql.NullString `db:"external_message_id"`
	DeletedAt         sql.NullTime   `db:"deleted_at"`
}

func (rcv *Message) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MessageFieldID,
		MessageFieldChatID,
		MessageFieldAuthorID,
		MessageFieldMessage,
		MessageFieldCreatedAt,
		MessageFieldUpdatedAt,
		MessageFieldDeletedAt,
		MessageFieldAttachments,
		MessageFieldExternalMessageID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ChatID,
		&rcv.AuthorID,
		&rcv.Message,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.Attachments,
		&rcv.ExternalMessageID,
	}

	return
}

func (*Message) TableName() string {
	return "messages"
}
