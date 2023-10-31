package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	StageStudentFieldID        = "id"
	StageStudentFieldStudentID = "student_id"
	StageStudentFieldCourseID  = "course_id"
	StageStudentFieldStageID   = "stage_id"
	StageStudentFieldStatus    = "status"
	StageStudentFieldCreatedBy = "created_by"
	StageStudentFieldUpdatedBy = "updated_by"
	StageStudentFieldCreatedAt = "created_at"
	StageStudentFieldUpdatedAt = "updated_at"
	StageStudentFieldDeletedAt = "deleted_at"
)

type StageStudent struct {
	ID        string                  `db:"id"`
	StudentID string                  `db:"student_id"`
	CourseID  string                  `db:"course_id"`
	StageID   string                  `db:"stage_id"`
	Status    enum.StageStudentStatus `db:"status"`
	CreatedBy string                  `db:"created_by"`
	UpdatedBy string                  `db:"updated_by"`
	CreatedAt time.Time               `db:"created_at"`
	UpdatedAt time.Time               `db:"updated_at"`
	DeletedAt sql.NullTime            `db:"deleted_at"`
}

func (rcv *StageStudent) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		StageStudentFieldID,
		StageStudentFieldStudentID,
		StageStudentFieldCourseID,
		StageStudentFieldStageID,
		StageStudentFieldStatus,
		StageStudentFieldCreatedBy,
		StageStudentFieldUpdatedBy,
		StageStudentFieldCreatedAt,
		StageStudentFieldUpdatedAt,
		StageStudentFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.StudentID,
		&rcv.CourseID,
		&rcv.StageID,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*StageStudent) TableName() string {
	return "stage_students"
}
