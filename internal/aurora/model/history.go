package model

import (
	"database/sql"
	"time"
)

const (
	HistoryFieldID         = "id"
	HistoryFieldTable      = "table_name"
	HistoryFieldRowID      = "row_id"
	HistoryFieldColumnName = "column_name"
	HistoryFieldOldValue   = "old_value"
	HistoryFieldNewValue   = "new_value"
	HistoryFieldCreatedAt  = "created_at"
)

type History struct {
	ID         int64                  `db:"id"`
	Table      string                 `db:"table_name"`
	RowID      string                 `db:"row_id"`
	ColumnName sql.NullString         `db:"column_name"`
	OldValue   map[string]interface{} `db:"old_value"`
	NewValue   map[string]interface{} `db:"new_value"`
	CreatedAt  time.Time              `db:"created_at"`
}

func (rcv *History) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		HistoryFieldID,
		HistoryFieldTable,
		HistoryFieldRowID,
		HistoryFieldColumnName,
		HistoryFieldOldValue,
		HistoryFieldNewValue,
		HistoryFieldCreatedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Table,
		&rcv.RowID,
		&rcv.ColumnName,
		&rcv.OldValue,
		&rcv.NewValue,
		&rcv.CreatedAt,
	}

	return
}

func (*History) TableName() string {
	return "history"
}
