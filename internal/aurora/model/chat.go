package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ChatFieldID                = "id"
	ChatFieldTitle             = "title"
	ChatFieldColor             = "color"
	ChatFieldExtranet          = "extranet"
	ChatFieldAuthorID          = "author_id"
	ChatFieldEntityType        = "entity_type"
	ChatFieldEntityID          = "entity_id"
	ChatFieldMessageCount      = "message_count"
	ChatFieldUserCount         = "user_count"
	ChatFieldPreMessageID      = "pre_message_id"
	ChatFieldLastMessageID     = "last_message_id"
	ChatFieldLastMessageStatus = "last_message_status"
	ChatFieldCreatedAt         = "created_at"
	ChatFieldUpdatedAt         = "updated_at"
	ChatFieldDeletedAt         = "deleted_at"
)

type Chat struct {
	ID                string              `db:"id"`
	Title             sql.NullString      `db:"title"`
	Color             sql.NullString      `db:"color"`
	Extranet          bool                `db:"extranet"`
	AuthorID          sql.NullString      `db:"author_id"`
	EntityType        enum.ChatEntityType `db:"entity_type"`
	EntityID          sql.NullString      `db:"entity_id"`
	MessageCount      int64               `db:"message_count"`
	UserCount         int64               `db:"user_count"`
	PreMessageID      sql.NullString      `db:"pre_message_id"`
	LastMessageID     sql.NullString      `db:"last_message_id"`
	LastMessageStatus enum.MessageStatus  `db:"last_message_status"`
	CreatedAt         time.Time           `db:"created_at"`
	UpdatedAt         time.Time           `db:"updated_at"`
	DeletedAt         sql.NullTime        `db:"deleted_at"`
}

func (rcv *Chat) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ChatFieldID,
		ChatFieldTitle,
		ChatFieldColor,
		ChatFieldExtranet,
		ChatFieldAuthorID,
		ChatFieldEntityType,
		ChatFieldEntityID,
		ChatFieldMessageCount,
		ChatFieldUserCount,
		ChatFieldPreMessageID,
		ChatFieldLastMessageID,
		ChatFieldLastMessageStatus,
		ChatFieldCreatedAt,
		ChatFieldUpdatedAt,
		ChatFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.Color,
		&rcv.Extranet,
		&rcv.AuthorID,
		&rcv.EntityType,
		&rcv.EntityID,
		&rcv.MessageCount,
		&rcv.UserCount,
		&rcv.PreMessageID,
		&rcv.LastMessageID,
		&rcv.LastMessageStatus,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Chat) TableName() string {
	return "chats"
}
