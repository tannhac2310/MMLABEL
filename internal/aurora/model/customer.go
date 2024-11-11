package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CustomerFieldID                 = "id"
	CustomerFieldName               = "name"
	CustomerFieldTax                = "tax"
	CustomerFieldCode               = "code"
	CustomerFieldCountry            = "country"
	CustomerFieldProvince           = "province"
	CustomerFieldAddress            = "address"
	CustomerFieldFax                = "fax"
	CustomerFieldCompanyWebsite     = "company_website"
	CustomerFieldCompanyPhone       = "company_phone"
	CustomerFieldCompanyEmail       = "company_email"
	CustomerFieldContactPersonName  = "contact_person_name"
	CustomerFieldContactPersonEmail = "contact_person_email"
	CustomerFieldContactPersonPhone = "contact_person_phone"
	CustomerFieldContactPersonRole  = "contact_person_role"
	CustomerFieldNote               = "note"
	CustomerFieldStatus             = "status"
	CustomerFieldData               = "data"
	CustomerFieldCreatedBy          = "created_by"
	CustomerFieldCreatedAt          = "created_at"
	CustomerFieldUpdatedAt          = "updated_at"
	CustomerFieldDeletedAt          = "deleted_at"
)

type Customer struct {
	ID                 string              `db:"id"`
	Name               string              `db:"name"`
	Tax                sql.NullString      `db:"tax"`
	Code               string              `db:"code"`
	Country            string              `db:"country"`
	Province           string              `db:"province"`
	Address            string              `db:"address"`
	Fax                sql.NullString      `db:"fax"`
	CompanyWebsite     sql.NullString      `db:"company_website"`
	CompanyPhone       sql.NullString      `db:"company_phone"`
	CompanyEmail       sql.NullString      `db:"company_email"`
	ContactPersonName  string              `db:"contact_person_name"`
	ContactPersonEmail string              `db:"contact_person_email"`
	ContactPersonPhone string              `db:"contact_person_phone"`
	ContactPersonRole  string              `db:"contact_person_role"`
	Note               sql.NullString      `db:"note"`
	Status             enum.CustomerStatus `db:"status"`
	Data               any                 `db:"data"`
	CreatedBy          string              `db:"created_by"`
	CreatedAt          time.Time           `db:"created_at"`
	UpdatedAt          time.Time           `db:"updated_at"`
	DeletedAt          sql.NullTime        `db:"deleted_at"`
}

func (rcv *Customer) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CustomerFieldID,
		CustomerFieldName,
		CustomerFieldTax,
		CustomerFieldCode,
		CustomerFieldCountry,
		CustomerFieldProvince,
		CustomerFieldAddress,
		CustomerFieldFax,
		CustomerFieldCompanyWebsite,
		CustomerFieldCompanyPhone,
		CustomerFieldCompanyEmail,
		CustomerFieldContactPersonName,
		CustomerFieldContactPersonEmail,
		CustomerFieldContactPersonPhone,
		CustomerFieldContactPersonRole,
		CustomerFieldNote,
		CustomerFieldStatus,
		CustomerFieldData,
		CustomerFieldCreatedBy,
		CustomerFieldCreatedAt,
		CustomerFieldUpdatedAt,
		CustomerFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Tax,
		&rcv.Code,
		&rcv.Country,
		&rcv.Province,
		&rcv.Address,
		&rcv.Fax,
		&rcv.CompanyWebsite,
		&rcv.CompanyPhone,
		&rcv.CompanyEmail,
		&rcv.ContactPersonName,
		&rcv.ContactPersonEmail,
		&rcv.ContactPersonPhone,
		&rcv.ContactPersonRole,
		&rcv.Note,
		&rcv.Status,
		&rcv.Data,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Customer) TableName() string {
	return "customers"
}
