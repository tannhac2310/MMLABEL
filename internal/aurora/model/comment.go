package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CommentFieldID         = "id"
	CommentFieldUserID     = "user_id"
	CommentFieldTargetID   = "target_id"
	CommentFieldTargetType = "target_type"
	CommentFieldContent    = "content"
	CommentFieldCreatedAt  = "created_at"
	CommentFieldUpdatedAt  = "updated_at"
	CommentFieldDeletedAt  = "deleted_at"
)

type Comment struct {
	ID         string             `db:"id"`
	UserID     string             `db:"user_id"`
	TargetID   string             `db:"target_id"`
	TargetType enum.CommentTarget `db:"target_type"`
	Content    string             `db:"content"`
	CreatedAt  time.Time          `db:"created_at"`
	UpdatedAt  time.Time          `db:"updated_at"`
	DeletedAt  sql.NullTime       `db:"deleted_at"`
}

func (rcv *Comment) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CommentFieldID,
		CommentFieldUserID,
		CommentFieldTargetID,
		CommentFieldTargetType,
		CommentFieldContent,
		CommentFieldCreatedAt,
		CommentFieldUpdatedAt,
		CommentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.UserID,
		&rcv.TargetID,
		&rcv.TargetType,
		&rcv.Content,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Comment) TableName() string {
	return "comments"
}
