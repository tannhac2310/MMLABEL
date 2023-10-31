package model

import (
	"time"

	"database/sql"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	OrderFieldID              = "id"
	OrderFieldInvoiceID       = "invoice_id"
	OrderFieldStudentID       = "student_id"
	OrderFieldStudentPhone    = "student_phone"
	OrderFieldStudentEmail    = "student_email"
	OrderFieldAddress         = "address"
	OrderFieldDetail          = "detail"
	OrderFieldPaymentMethod   = "payment_method"
	OrderFieldStatus          = "status"
	OrderFieldPrice           = "price"
	OrderFieldDiscount        = "discount"
	OrderFieldNote            = "note"
	OrderFieldCreatedBy       = "created_by"
	OrderFieldUpdatedBy       = "updated_by"
	OrderFieldCreatedAt       = "created_at"
	OrderFieldUpdatedAt       = "updated_at"
	OrderFieldDeletedAt       = "deleted_at"
	OrderFieldStudentFullName = "student_full_name"
)

type OrderDetail struct {
	CourseID     string            `json:"course_id"`
	StageID      string            `json:"stage_id"`
	Price        float32           `json:"price"`
	Discount     float32           `json:"discount"`
	DiscountType enum.DiscountType `json:"discount_type"`
	Note         string            `json:"note"`
}

type Order struct {
	ID              string             `db:"id"`
	InvoiceID       string             `db:"invoice_id"`
	StudentID       string             `db:"student_id"`
	StudentPhone    string             `db:"student_phone"`
	StudentEmail    string             `db:"student_email"`
	Address         string             `db:"address"`
	Detail          []*OrderDetail     `db:"detail"`
	PaymentMethod   enum.PaymentMethod `db:"payment_method"`
	Status          enum.OrderStatus   `db:"status"`
	Price           float32            `db:"price"`
	Discount        float32            `db:"discount"`
	Note            string             `db:"note"`
	CreatedBy       string             `db:"created_by"`
	UpdatedBy       string             `db:"updated_by"`
	CreatedAt       time.Time          `db:"created_at"`
	UpdatedAt       time.Time          `db:"updated_at"`
	DeletedAt       sql.NullTime       `db:"deleted_at"`
	StudentFullName string             `db:"student_full_name"`
}

func (rcv *Order) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		OrderFieldID,
		OrderFieldInvoiceID,
		OrderFieldStudentID,
		OrderFieldStudentPhone,
		OrderFieldStudentEmail,
		OrderFieldAddress,
		OrderFieldDetail,
		OrderFieldPaymentMethod,
		OrderFieldStatus,
		OrderFieldPrice,
		OrderFieldDiscount,
		OrderFieldNote,
		OrderFieldCreatedBy,
		OrderFieldUpdatedBy,
		OrderFieldCreatedAt,
		OrderFieldUpdatedAt,
		OrderFieldDeletedAt,
		OrderFieldStudentFullName,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.InvoiceID,
		&rcv.StudentID,
		&rcv.StudentPhone,
		&rcv.StudentEmail,
		&rcv.Address,
		&rcv.Detail,
		&rcv.PaymentMethod,
		&rcv.Status,
		&rcv.Price,
		&rcv.Discount,
		&rcv.Note,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.StudentFullName,
	}

	return
}

func (*Order) TableName() string {
	return "orders"
}
