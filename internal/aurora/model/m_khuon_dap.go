package model

import (
	"database/sql"
	"time"
)

const (
	MKhuonDapFieldID        = "id"
	MKhuonDapFieldData      = "data"
	MKhuonDapFieldCreatedBy = "created_by"
	MKhuonDapFieldCreatedAt = "created_at"
	MKhuonDapFieldUpdatedBy = "updated_by"
	MKhuonDapFieldUpdatedAt = "updated_at"
	MKhuonDapFieldDeletedAt = "deleted_at"
)

type MKhuonDap struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MKhuonDap) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MKhuonDapFieldID,
		MKhuonDapFieldData,
		MKhuonDapFieldCreatedBy,
		MKhuonDapFieldCreatedAt,
		MKhuonDapFieldUpdatedBy,
		MKhuonDapFieldUpdatedAt,
		MKhuonDapFieldDeletedAt,
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

func (*MKhuonDap) TableName() string {
	return "m_khuon_dap"
}
