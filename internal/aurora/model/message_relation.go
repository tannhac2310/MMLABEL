package model

import (
	"database/sql"
	"time"
)

const (
	MessageRelationFieldID        = "id"
	MessageRelationFieldChatID    = "chat_id"
	MessageRelationFieldUserID    = "user_id"
	MessageRelationFieldCreatedAt = "created_at"
	MessageRelationFieldUpdatedAt = "updated_at"
	MessageRelationFieldDeletedAt = "deleted_at"
)

type MessageRelation struct {
	ID        string       `db:"id"`
	ChatID    string       `db:"chat_id"`
	UserID    string       `db:"user_id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (rcv *MessageRelation) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MessageRelationFieldID,
		MessageRelationFieldChatID,
		MessageRelationFieldUserID,
		MessageRelationFieldCreatedAt,
		MessageRelationFieldUpdatedAt,
		MessageRelationFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ChatID,
		&rcv.UserID,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*MessageRelation) TableName() string {
	return "message_relations"
}
