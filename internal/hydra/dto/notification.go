package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type PushNotificationRequest struct {
	UserIDs        []string               `json:"userIds"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	Type           enum.NotificationType  `json:"type"`
	AdditionalData map[string]interface{} `json:"additionalData"`
}

type PushNotificationResponse struct {
}

type UserNotification struct {
	ID        string                  `json:"id"`
	UserID    string                  `json:"userId"`
	Title     string                  `json:"title"`
	Content   string                  `json:"content"`
	Type      enum.NotificationType   `json:"type"`
	Data      map[string]interface{}  `json:"data"`
	Status    enum.NotificationStatus `json:"status"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}
type UserNotificationFilter struct {
	IDs    []string                `json:"ids"`
	Type   enum.NotificationType   `json:"type"`
	Status enum.NotificationStatus `json:"status"`
}
type FindUserNotificationsRequest struct {
	Filter *UserNotificationFilter `json:"filter"`
	Paging *commondto.Paging       `json:"paging" binding:"required"`
}
type FindNotificationsResponse struct {
	Notifications []*UserNotification `json:"notifications"`
	NextPage      *commondto.Paging   `json:"nextPage"`
}

type CountUserNotificationRequests struct {
	Type   enum.NotificationType   `json:"type"`
	Status enum.NotificationStatus `json:"status"`
}

type CountUserNotificationResponse struct {
	TotalCount int64 `json:"totalCount"`
}

type MarkSeenNotificationsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}
type MarkSeenNotificationsResponse struct {
}

type MakeReadNotificationsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}
type MakeReadNotificationsResponse struct {
}
