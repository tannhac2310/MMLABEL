package model

import (
	"database/sql"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

const (
	InkFieldID             = "id"
	InkFieldImportID       = "import_id"
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
	InkFieldCreatedBy      = "created_by"
	InkFieldUpdatedBy      = "updated_by"
	InkFieldCreatedAt      = "created_at"
	InkFieldUpdatedAt      = "updated_at"
	InkFieldDeletedAt      = "deleted_at"
)

type Ink struct {
	ID             string                     `db:"id"`
	ImportID       sql.NullString             `db:"import_id"`
	Name           string                     `db:"name"`
	Code           string                     `db:"code"`
	ProductCodes   []string                   `db:"product_codes"`
	Position       string                     `db:"position"`
	Location       string                     `db:"location"`
	Manufacturer   string                     `db:"manufacturer"`
	ColorDetail    map[string]interface{}     `db:"color_detail"`
	Quantity       float64                    `db:"quantity"`
	ExpirationDate string                     `db:"expiration_date"` // DD-MM-YYYY
	Description    sql.NullString             `db:"description"`
	Data           map[string]interface{}     `db:"data"`
	Status         enum.InventoryCommonStatus `db:"status"`
	CreatedBy      string                     `db:"created_by"`
	UpdatedBy      string                     `db:"updated_by"`
	CreatedAt      time.Time                  `db:"created_at"`
	UpdatedAt      time.Time                  `db:"updated_at"`
	DeletedAt      sql.NullTime               `db:"deleted_at"`
}

func (rcv *Ink) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		InkFieldID,
		InkFieldImportID,
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
		InkFieldCreatedBy,
		InkFieldUpdatedBy,
		InkFieldCreatedAt,
		InkFieldUpdatedAt,
		InkFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ImportID,
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
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Ink) TableName() string {
	return "ink"
}
