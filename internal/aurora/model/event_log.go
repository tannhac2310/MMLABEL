package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	EventLogFieldID          = "id"
	EventLogFieldDeviceID    = "device_id"
	EventLogFieldStageID     = "stage_id"
	EventLogFieldStageStatus = "stage_status"
	EventLogFieldQuantity    = "quantity"
	EventLogFieldTimeSpend   = "time_spend"
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
	TimeSpend	int64							 `db:"time_spend"`
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
		EventLogFieldTimeSpend,
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
		&rcv.TimeSpend,
		&rcv.Msg,
		&rcv.Date,
		&rcv.CreatedAt,
	}

	return
}

func (*EventLog) TableName() string {
	return "event_logs"
}
