package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkFieldID             = "id"
	InkFieldName           = "name"
	InkFieldCode           = "code"
	InkFieldProductCodes   = "product_codes"
	InkFieldPosition       = "position"
	InkFieldLocation       = "location"
	InkFieldManufacturer   = "manufacturer"
	InkFieldColorDetail    = "color_detail"
	InkFieldQuantity       = "quantity"
	InkFieldExpirationDate = "expiration_date"
	InkFieldDescription    = "description"
	InkFieldData           = "data"
	InkFieldStatus         = "status"
	InkFieldCreatedAt      = "created_at"
	InkFieldUpdatedAt      = "updated_at"
	InkFieldDeletedAt      = "deleted_at"
)

type Ink struct {
	ID             string                           `db:"id"`
	Name           string                           `db:"name"`
	Code           string                           `db:"code"`
	ProductCodes   string                           `db:"product_codes"`
	Position       string                           `db:"position"`
	Location       string                           `db:"location"`
	Manufacturer   string                           `db:"manufacturer"`
	ColorDetail    map[string]interface{}           `db:"color_detail"`
	Quantity       float64                          `db:"quantity"`
	ExpirationDate time.Time                        `db:"expiration_date"`
	Description    sql.NullString                   `db:"description"`
	Data           map[string]interface{}           `db:"data"`
	Status         enum.InventoryCommonStatusStatus `db:"status"`
	CreatedAt      time.Time                        `db:"created_at"`
	UpdatedAt      time.Time                        `db:"updated_at"`
	DeletedAt      sql.NullTime                     `db:"deleted_at"`
}

func (rcv *Ink) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkFieldID,
		InkFieldName,
		InkFieldCode,
		InkFieldProductCodes,
		InkFieldPosition,
		InkFieldLocation,
		InkFieldManufacturer,
		InkFieldColorDetail,
		InkFieldQuantity,
		InkFieldExpirationDate,
		InkFieldDescription,
		InkFieldData,
		InkFieldStatus,
		InkFieldCreatedAt,
		InkFieldUpdatedAt,
		InkFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
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
		&rcv.Status,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Ink) TableName() string {
	return "ink"
}
