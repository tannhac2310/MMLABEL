package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	SMSHistoryFieldID        = "id"
	SMSHistoryFieldFrom      = "\"from\""
	SMSHistoryFieldTo        = "\"to\""
	SMSHistoryFieldContent   = "content"
	SMSHistoryFieldProvider  = "provider"
	SMSHistoryFieldRequest   = "request_data"
	SMSHistoryFieldResponse  = "response_data"
	SMSHistoryFieldCreatedAt = "created_at"
	SMSHistoryFieldUpdatedAt = "updated_at"
	SMSHistoryFieldDeletedAt = "deleted_at"
)

type SMSHistory struct {
	ID        string           `db:"id"`
	From      string           `db:"from"`
	To        string           `db:"to"`
	Content   string           `db:"content"`
	Provider  enum.SMSProvider `db:"provider"`
	Request   interface{}      `db:"request_data"`
	Response  interface{}      `db:"response_data"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
	DeletedAt sql.NullTime     `db:"deleted_at"`
}

func (*SMSHistory) TableName() string {
	return "sms_histories"
}

func (s *SMSHistory) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		SMSHistoryFieldID,
		SMSHistoryFieldFrom,
		SMSHistoryFieldTo,
		SMSHistoryFieldContent,
		SMSHistoryFieldProvider,
		SMSHistoryFieldRequest,
		SMSHistoryFieldResponse,
		SMSHistoryFieldCreatedAt,
		SMSHistoryFieldUpdatedAt,
		SMSHistoryFieldDeletedAt,
	}

	values = []interface{}{
		&s.ID,
		&s.From,
		&s.To,
		&s.Content,
		&s.Provider,
		&s.Request,
		&s.Response,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.DeletedAt,
	}

	return
}
