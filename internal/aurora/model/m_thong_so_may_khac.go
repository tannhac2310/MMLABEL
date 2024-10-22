package model

import (
	"database/sql"
	"time"
)

const (
	MThongSoMayKhacFieldID        = "id"
	MThongSoMayKhacFieldData      = "data"
	MThongSoMayKhacFieldCreatedBy = "created_by"
	MThongSoMayKhacFieldCreatedAt = "created_at"
	MThongSoMayKhacFieldUpdatedBy = "updated_by"
	MThongSoMayKhacFieldUpdatedAt = "updated_at"
	MThongSoMayKhacFieldDeletedAt = "deleted_at"
)

type MThongSoMayKhac struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MThongSoMayKhac) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MThongSoMayKhacFieldID,
		MThongSoMayKhacFieldData,
		MThongSoMayKhacFieldCreatedBy,
		MThongSoMayKhacFieldCreatedAt,
		MThongSoMayKhacFieldUpdatedBy,
		MThongSoMayKhacFieldUpdatedAt,
		MThongSoMayKhacFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Data,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*MThongSoMayKhac) TableName() string {
	return "m_thong_so_may_khac"
}
