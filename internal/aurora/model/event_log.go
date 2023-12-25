package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	EventLogFieldID          = "id"
	EventLogFieldDeviceID    = "device_id"
	EventLogFieldStageID     = "stage_id"
	EventLogFieldStageStatus = "stage_status"
	EventLogFieldQuantity    = "quantity"
	EventLogFieldMsg         = "msg"
	EventLogFieldDate        = "date"
	EventLogFieldCreatedAt   = "created_at"
)

type EventLog struct {
	ID          int64                            `db:"id"`
	DeviceID    string                           `db:"device_id"`
	StageID     sql.NullString                   `db:"stage_id"`
	StageStatus *enum.ProductionOrderStageStatus `db:"stage_status"`
	Quantity    float64                          `db:"quantity"`
	Msg         sql.NullString                   `db:"msg"`
	Date        sql.NullString                   `db:"date"`
	CreatedAt   time.Time                        `db:"created_at"`
}

func (rcv *EventLog) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		EventLogFieldID,
		EventLogFieldDeviceID,
		EventLogFieldStageID,
		EventLogFieldStageStatus,
		EventLogFieldQuantity,
		EventLogFieldMsg,
		EventLogFieldDate,
		EventLogFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.DeviceID,
		&rcv.StageID,
		&rcv.StageStatus,
		&rcv.Quantity,
		&rcv.Msg,
		&rcv.Date,
		&rcv.CreatedAt,
	}

	return
}

func (*EventLog) TableName() string {
	return "event_logs"
}
