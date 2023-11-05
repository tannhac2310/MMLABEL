package model

import (
	"database/sql"
	"time"
)

const (
	ProductionOrderStateDeviceEmployeeAssignmentFieldID                    = "id"
	ProductionOrderStateDeviceEmployeeAssignmentFieldPosDeviceAssignmentID = "pos_device_assignment_id"
	ProductionOrderStateDeviceEmployeeAssignmentFieldEmployeeID            = "employee_id"
	ProductionOrderStateDeviceEmployeeAssignmentFieldNote                  = "note"
	ProductionOrderStateDeviceEmployeeAssignmentFieldCreatedAt             = "created_at"
	ProductionOrderStateDeviceEmployeeAssignmentFieldUpdatedAt             = "updated_at"
	ProductionOrderStateDeviceEmployeeAssignmentFieldDeletedAt             = "deleted_at"
)

type ProductionOrderStateDeviceEmployeeAssignment struct {
	ID                    string         `db:"id"`
	PosDeviceAssignmentID sql.NullString `db:"pos_device_assignment_id"`
	EmployeeID            sql.NullString `db:"employee_id"`
	Note                  sql.NullString `db:"note"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
	DeletedAt             sql.NullTime   `db:"deleted_at"`
}

func (rcv *ProductionOrderStateDeviceEmployeeAssignment) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderStateDeviceEmployeeAssignmentFieldID,
		ProductionOrderStateDeviceEmployeeAssignmentFieldPosDeviceAssignmentID,
		ProductionOrderStateDeviceEmployeeAssignmentFieldEmployeeID,
		ProductionOrderStateDeviceEmployeeAssignmentFieldNote,
		ProductionOrderStateDeviceEmployeeAssignmentFieldCreatedAt,
		ProductionOrderStateDeviceEmployeeAssignmentFieldUpdatedAt,
		ProductionOrderStateDeviceEmployeeAssignmentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.PosDeviceAssignmentID,
		&rcv.EmployeeID,
		&rcv.Note,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*ProductionOrderStateDeviceEmployeeAssignment) TableName() string {
	return "production_order_state_device_employee_assignments"
}
