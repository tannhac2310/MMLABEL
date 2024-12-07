package model

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type DevicesError struct {
	DeviceID      string  `db:"device_id"`
	Quantity      int64   `db:"quantity"`
	QuantityError int64   `db:"quantity_error"`
	ErrorRate     float64 `db:"-"`
}

type DevicesProgressHistory struct {
	DeviceID      string                                `db:"device_id"`
	ProcessStatus enum.ProductionOrderStageDeviceStatus `db:"process_status"`
	CreatedAt     time.Time                             `db:"created_at"`
}

type DevicesErrorByStage struct {
	StageID       string `db:"stage_id"`
	QuantityError int64  `db:"quantity_error"`
}

type QuantityByDate struct {
	Day              int16 `db:"day"`
	QuantityComplete int64 `db:"quantity_complete"`
}

type QuantityDeliveryByDate struct {
	Day              int16 `db:"day"`
	QuantityDelivery int64 `db:"quantity_delivery"`
}
