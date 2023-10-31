package model

import (
	"database/sql"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	AttendanceDetailFieldID           = "id"
	AttendanceDetailFieldStudentID    = "student_id"
	AttendanceDetailFieldAttendanceID = "attendance_id"
	AttendanceDetailFieldStatus       = "status"
	AttendanceDetailFieldPoint        = "point"
	AttendanceDetailFieldDeletedAt    = "deleted_at"
)

type AttendanceDetail struct {
	ID           string                      `db:"id"`
	StudentID    string                      `db:"student_id"`
	AttendanceID string                      `db:"attendance_id"`
	Status       enum.AttendanceDetailStatus `db:"status"`
	Note         string                      `db:"note"`
	Point        float64                     `db:"point"`
	DeletedAt    sql.NullTime                `db:"deleted_at"`
}

func (rcv *AttendanceDetail) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		AttendanceDetailFieldID,
		AttendanceDetailFieldStudentID,
		AttendanceDetailFieldAttendanceID,
		AttendanceDetailFieldStatus,
		AttendanceDetailFieldPoint,
		AttendanceDetailFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.StudentID,
		&rcv.AttendanceID,
		&rcv.Status,
		&rcv.Point,
		&rcv.DeletedAt,
	}

	return
}

func (*AttendanceDetail) TableName() string {
	return "attendance_details"
}
