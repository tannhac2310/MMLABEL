package model

import (
	"database/sql"
	"time"
)

const (
	MKhuonBeFieldID        = "id"
	MKhuonBeFieldData      = "data"
	MKhuonBeFieldCreatedBy = "created_by"
	MKhuonBeFieldCreatedAt = "created_at"
	MKhuonBeFieldUpdatedBy = "updated_by"
	MKhuonBeFieldUpdatedAt = "updated_at"
	MKhuonBeFieldDeletedAt = "deleted_at"
)

type MKhuonBe struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MKhuonBe) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MKhuonBeFieldID,
		MKhuonBeFieldData,
		MKhuonBeFieldCreatedBy,
		MKhuonBeFieldCreatedAt,
		MKhuonBeFieldUpdatedBy,
		MKhuonBeFieldUpdatedAt,
		MKhuonBeFieldDeletedAt,
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

func (*MKhuonBe) TableName() string {
	return "m_khuon_be"
}
