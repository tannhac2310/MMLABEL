package nats

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type PushNotifyEvent struct {
	UserIDs        []string
	Title          string
	Content        string
	Type           enum.NotificationType
	AdditionalData map[string]interface{}
	SaveToDB       bool
	TimeToLiveSec  *uint
	Sound          string
	PushFCM        bool
}

type Coordinates struct {
	Latitude  string
	Longitude string
}
type PayloadAttachment struct {
	ID          string
	Thumbnail   string
	URL         string
	Description string
	Coordinates *Coordinates
	Size        string
	Name        string
	Checksum    string
	Type        string
}
type Attachment struct {
	Payload *PayloadAttachment
	Type    string
}
type ZaloMessageImageAttachments struct {
	Attachments []*Attachment
}
type ZaloSender struct {
	ID string
}
type ZaloMessageText struct {
	Text  string
	MsgID string
}
type ZaloEventRequest struct {
	AppID       string
	Sender      *ZaloSender
	UserIDByApp string
	Recipient   struct {
		ID string
	}
	Message     *ZaloMessageText
	Attachments []*Attachment
	EventName   string
	Timestamp   time.Time
}

type ZaloEventMessage struct {
	AppID              string
	ZaloEventSignature string
	EventName          string
	EventData          *ZaloEventRequest
}
