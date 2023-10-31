package dto



type ZaloSender struct {
	ID string `json:"id"`
}
type ZaloMessageText struct {
	Text  string `json:"text"`
	MsgID string `json:"msg_id"`
}
type ZaloEventRequest struct {
	AppID       string           `json:"app_id"`
	Sender      *ZaloSender      `json:"sender"`
	UserIDByApp string           `json:"user_id_by_app"`
	Message     *ZaloMessageText `json:"message"`
	Attachments []*Attachment    `json:"attachments"`
	Recipient   struct {
		ID string `json:"id"`
	} `json:"recipient"`
	EventName string `json:"event_name"`
	Timestamp string `json:"timestamp"`
}
type ZaloEventResponse struct {
}


