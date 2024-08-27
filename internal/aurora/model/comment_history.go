package model

import "time"

const (
	CommentHistoryFieldID        = "id"
	CommentHistoryFieldCommentID = "comment_id"
	CommentHistoryFieldContent   = "content"
	CommentHistoryFieldCreatedAt = "created_at"
)

type CommentHistory struct {
	ID        string    `db:"id"`
	CommentID string    `db:"comment_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func (rcv *CommentHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CommentHistoryFieldID,
		CommentHistoryFieldCommentID,
		CommentHistoryFieldContent,
		CommentHistoryFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.CommentID,
		&rcv.Content,
		&rcv.CreatedAt,
	}

	return
}

func (*CommentHistory) TableName() string {
	return "comment_histories"
}
