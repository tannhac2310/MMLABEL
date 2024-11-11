package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	MasterDataFieldID          = "id"
	MasterDataFieldType        = "type"
	MasterDataFieldName        = "name"
	MasterDataFieldCode        = "code"
	MasterDataFieldDescription = "description"
	MasterDataFieldStatus      = "status"
	MasterDataFieldCreatedAt   = "created_at"
	MasterDataFieldUpdatedAt   = "updated_at"
	MasterDataFieldCreatedBy   = "created_by"
	MasterDataFieldUpdatedBy   = "updated_by"
	MasterDataFieldDeletedAt   = "deleted_at"
)

type MasterData struct {
	ID          string                `db:"id"`
	Type        enum.MasterDataType   `db:"type"`
	Name        string                `db:"name"`
	Code        string                `db:"code"`
	Description string                `db:"description"`
	Status      enum.MasterDataStatus `db:"status"`
	CreatedAt   time.Time             `db:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at"`
	CreatedBy   string                `db:"created_by"`
	UpdatedBy   string                `db:"updated_by"`
	DeletedAt   sql.NullTime          `db:"deleted_at"`
}

func (rcv *MasterData) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MasterDataFieldID,
		MasterDataFieldType,
		MasterDataFieldName,
		MasterDataFieldCode,
		MasterDataFieldDescription,
		MasterDataFieldStatus,
		MasterDataFieldCreatedAt,
		MasterDataFieldUpdatedAt,
		MasterDataFieldCreatedBy,
		MasterDataFieldUpdatedBy,
		MasterDataFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Type,
		&rcv.Name,
		&rcv.Code,
		&rcv.Description,
		&rcv.Status,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.DeletedAt,
	}

	return
}

func (*MasterData) TableName() string {
	return "master_data"
}
