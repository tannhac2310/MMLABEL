package ws

const (
	// WebsocketAction
	WebsocketActionPing      = "ping"
	WebsocketActionBroadcast = "broadcast"

	// WebsocketEvent
	WebsocketEventResponse              = "response"
	WebsocketEventBroadcast             = "broadcast"
	WebsocketEventIoTDeviceStateChange  = "iot_device_state_change"
	WebsocketEventPondModeChange        = "pond_mode_change"
	WebsocketEventEdgeDeviceStateChange = "edge_device_state_change"

	// status
	StatusOk   = "OK"
	StatusFail = "FAIL"
)
