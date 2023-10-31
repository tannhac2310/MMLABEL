package model

import (
	"database/sql"
	"time"
)

const (
	BannerFieldID           = "id"
	BannerFieldName         = "name"
	BannerFieldLink         = "link"
	BannerFieldDisplayOrder = "display_order"
	BannerFieldPhotoURL     = "photo_url"
	BannerFieldCreatedBy    = "created_by"
	BannerFieldCreatedAt    = "created_at"
	BannerFieldUpdatedAt    = "updated_at"
	BannerFieldDeletedAt    = "deleted_at"
)

type Banner struct {
	ID           string       `db:"id"`
	Link         string       `db:"link"`
	DisplayOrder int8         `db:"display_order"`
	Name         string       `db:"name"`
	PhotoURL     string       `db:"photo_url"`
	CreatedBy    string       `db:"created_by"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

func (b *Banner) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		BannerFieldID,
		BannerFieldName,
		BannerFieldLink,
		BannerFieldDisplayOrder,
		BannerFieldPhotoURL,
		BannerFieldCreatedBy,
		BannerFieldCreatedAt,
		BannerFieldUpdatedAt,
		BannerFieldDeletedAt,
	}

	values = []interface{}{
		&b.ID,
		&b.Name,
		&b.Link,
		&b.DisplayOrder,
		&b.PhotoURL,
		&b.CreatedBy,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	}
	return
}

func (*Banner) TableName() string {
	return "banners"
}
