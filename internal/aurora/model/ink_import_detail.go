package model

import (
	"database/sql"
	"time"
)

const (
	InkImportDetailFieldID             = "id"
	InkImportDetailFieldInkID          = "ink_id"
	InkImportDetailFieldInkImportID    = "ink_import_id"
	InkImportDetailFieldName           = "name"
	InkImportDetailFieldCode           = "code"
	InkImportDetailFieldProductCodes   = "product_codes"
	InkImportDetailFieldPosition       = "position"
	InkImportDetailFieldLocation       = "location"
	InkImportDetailFieldManufacturer   = "manufacturer"
	InkImportDetailFieldColorDetail    = "color_detail"
	InkImportDetailFieldQuantity       = "quantity"
	InkImportDetailFieldExpirationDate = "expiration_date"
	InkImportDetailFieldDescription    = "description"
	InkImportDetailFieldData           = "data"
	InkImportDetailFieldCreatedAt      = "created_at"
	InkImportDetailFieldUpdatedAt      = "updated_at"
	InkImportDetailFieldDeletedAt      = "deleted_at"
)

type InkImportDetail struct {
	ID             string                 `db:"id"`
	InkID          string                 `db:"ink_id"`
	InkImportID    string                 `db:"ink_import_id"`
	Name           string                 `db:"name"`
	Code           string                 `db:"code"`
	ProductCodes   []string               `db:"product_codes"`
	Position       string                 `db:"position"`
	Location       string                 `db:"location"`
	Manufacturer   string                 `db:"manufacturer"`
	ColorDetail    map[string]interface{} `db:"color_detail"`
	Quantity       float64                `db:"quantity"`
	ExpirationDate string                 `db:"expiration_date"` // DD-MM-YYYY
	Description    sql.NullString         `db:"description"`
	Data           map[string]interface{} `db:"data"`
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at"`
	DeletedAt      sql.NullTime           `db:"deleted_at"`
}

func (rcv *InkImportDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkImportDetailFieldID,
		InkImportDetailFieldInkID,
		InkImportDetailFieldInkImportID,
		InkImportDetailFieldName,
		InkImportDetailFieldCode,
		InkImportDetailFieldProductCodes,
		InkImportDetailFieldPosition,
		InkImportDetailFieldLocation,
		InkImportDetailFieldManufacturer,
		InkImportDetailFieldColorDetail,
		InkImportDetailFieldQuantity,
		InkImportDetailFieldExpirationDate,
		InkImportDetailFieldDescription,
		InkImportDetailFieldData,
		InkImportDetailFieldCreatedAt,
		InkImportDetailFieldUpdatedAt,
		InkImportDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.InkID,
		&rcv.InkImportID,
		&rcv.Name,
		&rcv.Code,
		&rcv.ProductCodes,
		&rcv.Position,
		&rcv.Location,
		&rcv.Manufacturer,
		&rcv.ColorDetail,
		&rcv.Quantity,
		&rcv.ExpirationDate,
		&rcv.Description,
		&rcv.Data,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*InkImportDetail) TableName() string {
	return "ink_import_detail"
}
