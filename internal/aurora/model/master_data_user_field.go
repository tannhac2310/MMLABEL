package model

import (
	"database/sql"
	"time"
)

const (
	MasterDataUserFieldFieldID           = "id"
	MasterDataUserFieldFieldMasterDataID = "master_data_id"
	MasterDataUserFieldFieldFieldName    = "field_name"
	MasterDataUserFieldFieldFieldValue   = "field_value"
	MasterDataUserFieldFieldData         = "data"
	MasterDataUserFieldFieldCreatedAt    = "created_at"
	MasterDataUserFieldFieldUpdatedAt    = "updated_at"
	MasterDataUserFieldFieldCreatedBy    = "created_by"
	MasterDataUserFieldFieldUpdatedBy    = "updated_by"
	MasterDataUserFieldFieldDeletedAt    = "deleted_at"
)

type MasterDataUserField struct {
	ID           string       `db:"id"`
	MasterDataID string       `db:"master_data_id"`
	FieldName    string       `db:"field_name"`
	FieldValue   string       `db:"field_value"`
	Data         any          `db:"data"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	CreatedBy    string       `db:"created_by"`
	UpdatedBy    string       `db:"updated_by"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

func (rcv *MasterDataUserField) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MasterDataUserFieldFieldID,
		MasterDataUserFieldFieldMasterDataID,
		MasterDataUserFieldFieldFieldName,
		MasterDataUserFieldFieldFieldValue,
		MasterDataUserFieldFieldData,
		MasterDataUserFieldFieldCreatedAt,
		MasterDataUserFieldFieldUpdatedAt,
		MasterDataUserFieldFieldCreatedBy,
		MasterDataUserFieldFieldUpdatedBy,
		MasterDataUserFieldFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.MasterDataID,
		&rcv.FieldName,
		&rcv.FieldValue,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.DeletedAt,
	}

	return
}

func (*MasterDataUserField) TableName() string {
	return "master_data_user_field"
}
