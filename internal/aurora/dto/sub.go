package dto

type ZaloSender struct {
	ID string `json:"id"`
}
type ZaloMessageText struct {
	Text  string `json:"text"`
	MsgID string `json:"msg_id"`
}
type ZaloEventResponse struct {
}
