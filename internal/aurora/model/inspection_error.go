package model

import (
	"time"
)

const (
	InspectionErrorFieldID               = "id"
	InspectionErrorFieldDeviceID         = "device_id"
	InspectionErrorFieldDeviceName       = "device_name"
	InspectionErrorFieldInspectionFormID = "inspection_form_id"
	InspectionErrorFieldErrorType        = "error_type"
	InspectionErrorFieldQuantity         = "quantity"
	InspectionErrorFieldNhanVienThucHien = "nhan_vien_thuc_hien"
	InspectionErrorFieldNote             = "note"
	InspectionErrorFieldCreatedBy        = "created_by"
	InspectionErrorFieldUpdatedBy        = "updated_by"
	InspectionErrorFieldCreatedAt        = "created_at"
	InspectionErrorFieldUpdatedAt        = "updated_at"
)

type InspectionError struct {
	ID               string    `db:"id"`
	DeviceID         string    `db:"device_id"`
	DeviceName       string    `db:"device_name"`
	InspectionFormID string    `db:"inspection_form_id"`
	ErrorType        string    `db:"error_type"`
	Quantity         int64     `db:"quantity"`
	NhanVienThucHien string    `db:"nhan_vien_thuc_hien"`
	Note             string    `db:"note"`
	CreatedBy        string    `db:"created_by"`
	UpdatedBy        string    `db:"updated_by"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (rcv *InspectionError) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InspectionErrorFieldID,
		InspectionErrorFieldDeviceID,
		InspectionErrorFieldDeviceName,
		InspectionErrorFieldInspectionFormID,
		InspectionErrorFieldErrorType,
		InspectionErrorFieldQuantity,
		InspectionErrorFieldNhanVienThucHien,
		InspectionErrorFieldNote,
		InspectionErrorFieldCreatedBy,
		InspectionErrorFieldUpdatedBy,
		InspectionErrorFieldCreatedAt,
		InspectionErrorFieldUpdatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.DeviceID,
		&rcv.DeviceName,
		&rcv.InspectionFormID,
		&rcv.ErrorType,
		&rcv.Quantity,
		&rcv.NhanVienThucHien,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
	}

	return
}

func (*InspectionError) TableName() string {
	return "inspection_errors"
}
