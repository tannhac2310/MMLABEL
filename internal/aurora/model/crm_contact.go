package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CrmContactFieldID          = "id"
	CrmContactFieldName        = "name"
	CrmContactFieldAvatar      = "avatar"
	CrmContactFieldPhoneNumber = "phone_number"
	CrmContactFieldEmail       = "email"
	CrmContactFieldFacebook    = "facebook"
	CrmContactFieldStatus      = "status"
	CrmContactFieldAddress     = "address"
	CrmContactFieldSourceType  = "source_type"
	CrmContactFieldSourceID    = "source_id"
	CrmContactFieldExternalID  = "external_id"
	CrmContactFieldCreatedBy   = "created_by"
	CrmContactFieldUpdatedBy   = "updated_by"
	CrmContactFieldCreatedAt   = "created_at"
	CrmContactFieldUpdatedAt   = "updated_at"
	CrmContactFieldDeletedAt   = "deleted_at"
)

type CrmContact struct {
	ID          string                `db:"id"`
	Name        string                `db:"name"`
	Avatar      sql.NullString        `db:"avatar"`
	PhoneNumber sql.NullString        `db:"phone_number"`
	Email       sql.NullString        `db:"email"`
	Facebook    sql.NullString        `db:"facebook"`
	Status      enum.CommonStatus     `db:"status"`
	Address     sql.NullString        `db:"address"`
	SourceType  enum.CrmContactSource `db:"source_type"`
	SourceID    string                `db:"source_id"`
	ExternalID  sql.NullString        `db:"external_id"`
	CreatedBy   sql.NullString        `db:"created_by"`
	UpdatedBy   sql.NullString        `db:"updated_by"`
	CreatedAt   time.Time             `db:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at"`
	DeletedAt   sql.NullTime          `db:"deleted_at"`
}

func (rcv *CrmContact) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CrmContactFieldID,
		CrmContactFieldName,
		CrmContactFieldAvatar,
		CrmContactFieldPhoneNumber,
		CrmContactFieldEmail,
		CrmContactFieldFacebook,
		CrmContactFieldStatus,
		CrmContactFieldAddress,
		CrmContactFieldSourceType,
		CrmContactFieldSourceID,
		CrmContactFieldExternalID,
		CrmContactFieldCreatedBy,
		CrmContactFieldUpdatedBy,
		CrmContactFieldCreatedAt,
		CrmContactFieldUpdatedAt,
		CrmContactFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Avatar,
		&rcv.PhoneNumber,
		&rcv.Email,
		&rcv.Facebook,
		&rcv.Status,
		&rcv.Address,
		&rcv.SourceType,
		&rcv.SourceID,
		&rcv.ExternalID,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*CrmContact) TableName() string {
	return "crm_contacts"
}
