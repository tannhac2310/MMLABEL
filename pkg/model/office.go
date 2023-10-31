package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	OfficeFieldID          = "id"
	OfficeFieldName        = "name"
	OfficeFieldPhone       = "phone"
	OfficeFieldAddress     = "address"
	OfficeFieldProvinceID  = "province_id"
	OfficeFieldDistrictID  = "district_id"
	OfficeFieldStatus      = "status"
	OfficeFieldPhotoURL    = "photo_url"
	OfficeFieldDescription = "description"
	OfficeFieldCreatedBy   = "created_by"
	OfficeFieldCreatedAt   = "created_at"
	OfficeFieldUpdatedAt   = "updated_at"
	OfficeFieldDeletedAt   = "deleted_at"
)

type Office struct {
	ID          string            `db:"id"`
	Name        string            `db:"name"`
	Phone       string            `db:"phone"`
	Address     string            `db:"address"`
	ProvinceID  int64             `db:"province_id"`
	DistrictID  int64             `db:"district_id"`
	Status      enum.CommonStatus `db:"status"`
	PhotoURL    string            `db:"photo_url"`
	Description string            `db:"description"`
	CreatedBy   string            `db:"created_by"`
	CreatedAt   time.Time         `db:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"`
	DeletedAt   sql.NullTime      `db:"deleted_at"`
}

func (b *Office) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OfficeFieldID,
		OfficeFieldName,
		OfficeFieldPhone,
		OfficeFieldAddress,
		OfficeFieldProvinceID,
		OfficeFieldDistrictID,
		OfficeFieldStatus,
		OfficeFieldPhotoURL,
		OfficeFieldDescription,
		OfficeFieldCreatedBy,
		OfficeFieldCreatedAt,
		OfficeFieldUpdatedAt,
		OfficeFieldDeletedAt,
	}

	values = []interface{}{
		&b.ID,
		&b.Name,
		&b.Phone,
		&b.Address,
		&b.ProvinceID,
		&b.DistrictID,
		&b.Status,
		&b.PhotoURL,
		&b.Description,
		&b.CreatedBy,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	}
	return
}

func (*Office) TableName() string {
	return "offices"
}
