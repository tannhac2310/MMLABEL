package model

import (
	"database/sql"
	"time"
)

const (
	MPhimFieldID        = "id"
	MPhimFieldData      = "data"
	MPhimFieldCreatedBy = "created_by"
	MPhimFieldCreatedAt = "created_at"
	MPhimFieldUpdatedBy = "updated_by"
	MPhimFieldUpdatedAt = "updated_at"
	MPhimFieldDeletedAt = "deleted_at"
)

type MPhim struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MPhim) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MPhimFieldID,
		MPhimFieldData,
		MPhimFieldCreatedBy,
		MPhimFieldCreatedAt,
		MPhimFieldUpdatedBy,
		MPhimFieldUpdatedAt,
		MPhimFieldDeletedAt,
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

func (*MPhim) TableName() string {
	return "m_phim"
}
