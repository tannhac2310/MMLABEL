package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	StageFieldID             = "id"
	StageFieldParentID       = "parent_id"
	StageFieldDepartmentCode = "department_code"
	StageFieldName           = "name"
	StageFieldShortName      = "short_name"
	StageFieldCode           = "code"
	StageFieldSorting        = "sorting"
	StageFieldErrorCode      = "error_code"
	StageFieldData           = "data"
	StageFieldNote           = "note"
	StageFieldStatus         = "status"
	StageFieldCreatedBy      = "created_by"
	StageFieldCreatedAt      = "created_at"
	StageFieldUpdatedAt      = "updated_at"
	StageFieldDeletedAt      = "deleted_at"
)

type Stage struct {
	ID             string                 `db:"id"`
	ParentID       sql.NullString         `db:"parent_id"`
	DepartmentCode sql.NullString         `db:"department_code"`
	Name           string                 `db:"name"`
	ShortName      string                 `db:"short_name"`
	Code           string                 `db:"code"`
	Sorting        int16                  `db:"sorting"`
	ErrorCode      sql.NullString         `db:"error_code"`
	Data           map[string]interface{} `db:"data"`
	Note           sql.NullString         `db:"note"`
	Status         enum.StageStatus       `db:"status"`
	CreatedBy      string                 `db:"created_by"`
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at"`
	DeletedAt      sql.NullTime           `db:"deleted_at"`
}

func (rcv *Stage) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		StageFieldID,
		StageFieldParentID,
		StageFieldDepartmentCode,
		StageFieldName,
		StageFieldShortName,
		StageFieldCode,
		StageFieldSorting,
		StageFieldErrorCode,
		StageFieldData,
		StageFieldNote,
		StageFieldStatus,
		StageFieldCreatedBy,
		StageFieldCreatedAt,
		StageFieldUpdatedAt,
		StageFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ParentID,
		&rcv.DepartmentCode,
		&rcv.Name,
		&rcv.ShortName,
		&rcv.Code,
		&rcv.Sorting,
		&rcv.ErrorCode,
		&rcv.Data,
		&rcv.Note,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Stage) TableName() string {
	return "stages"
}
