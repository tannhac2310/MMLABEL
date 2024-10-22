package model

import (
	"database/sql"
	"time"
)

const (
	MKhungInFieldID        = "id"
	MKhungInFieldData      = "data"
	MKhungInFieldCreatedBy = "created_by"
	MKhungInFieldCreatedAt = "created_at"
	MKhungInFieldUpdatedBy = "updated_by"
	MKhungInFieldUpdatedAt = "updated_at"
	MKhungInFieldDeletedAt = "deleted_at"
)

type MKhungIn struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MKhungIn) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MKhungInFieldID,
		MKhungInFieldData,
		MKhungInFieldCreatedBy,
		MKhungInFieldCreatedAt,
		MKhungInFieldUpdatedBy,
		MKhungInFieldUpdatedAt,
		MKhungInFieldDeletedAt,
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

func (*MKhungIn) TableName() string {
	return "m_khung_in"
}
