package model

import (
	"database/sql"
	"time"
)

const (
	MNguyenVatLieuFieldID        = "id"
	MNguyenVatLieuFieldData      = "data"
	MNguyenVatLieuFieldCreatedBy = "created_by"
	MNguyenVatLieuFieldCreatedAt = "created_at"
	MNguyenVatLieuFieldUpdatedBy = "updated_by"
	MNguyenVatLieuFieldUpdatedAt = "updated_at"
	MNguyenVatLieuFieldDeletedAt = "deleted_at"
)

type MNguyenVatLieu struct {
	ID        string                 `db:"id"`
	Data      map[string]interface{} `db:"data"`
	CreatedBy string                 `db:"created_by"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedBy string                 `db:"updated_by"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt sql.NullTime           `db:"deleted_at"`
}

func (rcv *MNguyenVatLieu) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MNguyenVatLieuFieldID,
		MNguyenVatLieuFieldData,
		MNguyenVatLieuFieldCreatedBy,
		MNguyenVatLieuFieldCreatedAt,
		MNguyenVatLieuFieldUpdatedBy,
		MNguyenVatLieuFieldUpdatedAt,
		MNguyenVatLieuFieldDeletedAt,
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

func (*MNguyenVatLieu) TableName() string {
	return "m_nguyen_vat_lieu"
}
