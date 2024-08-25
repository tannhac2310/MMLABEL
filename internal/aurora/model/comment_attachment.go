package model

import (
	"database/sql"
	"time"
)

const (
	CommentAttachmentFieldID        = "id"
	CommentAttachmentFieldCommentID = "comment_id"
	CommentAttachmentFieldUrl       = "url"
	CommentAttachmentFieldCreatedAt = "created_at"
	CommentAttachmentFieldDeletedAt = "deleted_at"
)

type CommentAttachment struct {
	ID        string       `db:"id"`
	CommentID string       `db:"comment_id"`
	Url       string       `db:"url"`
	CreatedAt time.Time    `db:"created_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (rcv *CommentAttachment) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CommentAttachmentFieldID,
		CommentAttachmentFieldCommentID,
		CommentAttachmentFieldUrl,
		CommentAttachmentFieldCreatedAt,
		CommentAttachmentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.CommentID,
		&rcv.Url,
		&rcv.CreatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*CommentAttachment) TableName() string {
	return "comment_attachments"
}
