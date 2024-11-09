package model

import (
	"database/sql"
	"time"
)

const (
	MasterDataSelectionFieldID              = "id"
	MasterDataSelectionFieldSelectionGroup  = "selection_group"
	MasterDataSelectionFieldDisplayValue    = "display_value"
	MasterDataSelectionFieldValue           = "value"
	MasterDataSelectionFieldDescription     = "description"
	MasterDataSelectionFieldMultipleChoices = "multiple_choices"
	MasterDataSelectionFieldSortOrder       = "sort_order"
	MasterDataSelectionFieldStatus          = "status"
	MasterDataSelectionFieldCreatedAt       = "created_at"
	MasterDataSelectionFieldUpdatedAt       = "updated_at"
	MasterDataSelectionFieldCreatedBy       = "created_by"
	MasterDataSelectionFieldUpdatedBy       = "updated_by"
	MasterDataSelectionFieldDeletedAt       = "deleted_at"
)

type MasterDataSelection struct {
	ID              string         `db:"id"`
	SelectionGroup  string         `db:"selection_group"`
	DisplayValue    string         `db:"display_value"`
	Value           string         `db:"value"`
	Description     sql.NullString `db:"description"`
	MultipleChoices int16          `db:"multiple_choices"`
	SortOrder       int16          `db:"sort_order"`
	Status          int16          `db:"status"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
	CreatedBy       sql.NullString `db:"created_by"`
	UpdatedBy       sql.NullString `db:"updated_by"`
	DeletedAt       sql.NullTime   `db:"deleted_at"`
}

func (rcv *MasterDataSelection) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		MasterDataSelectionFieldID,
		MasterDataSelectionFieldSelectionGroup,
		MasterDataSelectionFieldDisplayValue,
		MasterDataSelectionFieldValue,
		MasterDataSelectionFieldDescription,
		MasterDataSelectionFieldMultipleChoices,
		MasterDataSelectionFieldSortOrder,
		MasterDataSelectionFieldStatus,
		MasterDataSelectionFieldCreatedAt,
		MasterDataSelectionFieldUpdatedAt,
		MasterDataSelectionFieldCreatedBy,
		MasterDataSelectionFieldUpdatedBy,
		MasterDataSelectionFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.SelectionGroup,
		&rcv.DisplayValue,
		&rcv.Value,
		&rcv.Description,
		&rcv.MultipleChoices,
		&rcv.SortOrder,
		&rcv.Status,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.DeletedAt,
	}

	return
}

func (*MasterDataSelection) TableName() string {
	return "master_data_selection"
}
