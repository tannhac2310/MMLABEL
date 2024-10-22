package model

import (
	"database/sql"
	"time"
)

const (
	MThongSoMayInFieldID        = "id"
	MThongSoMayInFieldData      = "data"
	MThongSoMayInFieldCreatedBy = "created_by"
	MThongSoMayInFieldCreatedAt = "created_at"
	MThongSoMayInFieldUpdatedBy = "updated_by"
	MThongSoMayInFieldUpdatedAt = "updated_at"
	MThongSoMayInFieldDeletedAt = "deleted_at"
)

type MThongSoMayIn struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MThongSoMayIn) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MThongSoMayInFieldID,
		MThongSoMayInFieldData,
		MThongSoMayInFieldCreatedBy,
		MThongSoMayInFieldCreatedAt,
		MThongSoMayInFieldUpdatedBy,
		MThongSoMayInFieldUpdatedAt,
		MThongSoMayInFieldDeletedAt,
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

func (*MThongSoMayIn) TableName() string {
	return "m_thong_so_may_in"
}
