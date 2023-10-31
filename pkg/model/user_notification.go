package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	UserNotificationFieldID        = "id"
	UserNotificationFieldUserID    = "user_id"
	UserNotificationFieldTitle     = "title"
	UserNotificationFieldContent   = "content"
	UserNotificationFieldType      = "type"
	UserNotificationFieldData      = "data"
	UserNotificationFieldStatus    = "status"
	UserNotificationFieldCreatedAt = "created_at"
	UserNotificationFieldUpdatedAt = "updated_at"
	UserNotificationFieldDeletedAt = "deleted_at"
)

type UserNotification struct {
	ID        string                  `db:"id"`
	UserID    string                  `db:"user_id"`
	Title     string                  `db:"title"`
	Content   string                  `db:"content"`
	Type      enum.NotificationType   `db:"type"`
	Data      map[string]interface{}  `db:"data"`
	Status    enum.NotificationStatus `db:"status"`
	CreatedAt time.Time               `db:"created_at"`
	UpdatedAt time.Time               `db:"updated_at"`
	DeletedAt sql.NullTime            `db:"deleted_at"`
}

func (*UserNotification) TableName() string {
	return "user_notifications"
}

func (u *UserNotification) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		UserNotificationFieldID,
		UserNotificationFieldUserID,
		UserNotificationFieldTitle,
		UserNotificationFieldContent,
		UserNotificationFieldType,
		UserNotificationFieldData,
		UserNotificationFieldStatus,
		UserNotificationFieldCreatedAt,
		UserNotificationFieldUpdatedAt,
		UserNotificationFieldDeletedAt,
	}

	values = []interface{}{
		&u.ID,
		&u.UserID,
		&u.Title,
		&u.Content,
		&u.Type,
		&u.Data,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	}

	return
}
