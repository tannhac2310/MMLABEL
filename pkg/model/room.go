package model

import (
	"database/sql"
	"time"
)

const (
	RoomFieldID          = "id"
	RoomFieldName        = "name"
	RoomFieldOfficeID    = "office_id"
	RoomFieldStatus      = "status"
	RoomFieldPhotoUrl    = "photo_url"
	RoomFieldDescription = "description"
	RoomFieldCreatedBy   = "created_by"
	RoomFieldUpdatedBy   = "updated_by"
	RoomFieldCreatedAt   = "created_at"
	RoomFieldUpdatedAt   = "updated_at"
	RoomFieldDeletedAt   = "deleted_at"
)

type Room struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	OfficeID    string         `db:"office_id"`
	Status      int16          `db:"status"`
	PhotoUrl    sql.NullString `db:"photo_url"`
	Description sql.NullString `db:"description"`
	CreatedBy   string         `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

func (rcv *Room) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		RoomFieldID,
		RoomFieldName,
		RoomFieldOfficeID,
		RoomFieldStatus,
		RoomFieldPhotoUrl,
		RoomFieldDescription,
		RoomFieldCreatedBy,
		RoomFieldUpdatedBy,
		RoomFieldCreatedAt,
		RoomFieldUpdatedAt,
		RoomFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.OfficeID,
		&rcv.Status,
		&rcv.PhotoUrl,
		&rcv.Description,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Room) TableName() string {
	return "rooms"
}
