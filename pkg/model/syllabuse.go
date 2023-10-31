package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	SyllabusFieldID          = "id"
	SyllabusFieldTitle       = "title"
	SyllabusFieldCode        = "code"
	SyllabusFieldTeacherID   = "teacher_id"
	SyllabusFieldCourseID    = "course_id"
	SyllabusFieldDescription = "description"
	SyllabusFieldStatus      = "status"
	SyllabusFieldCreatedBy   = "created_by"
	SyllabusFieldUpdatedBy   = "updated_by"
	SyllabusFieldCreatedAt   = "created_at"
	SyllabusFieldUpdatedAt   = "updated_at"
	SyllabusFieldDeletedAt   = "deleted_at"
)

type Syllabus struct {
	ID          string            `db:"id"`
	Title       string            `db:"title"`
	Code        string            `db:"code"`
	TeacherID   string            `db:"teacher_id"`
	CourseID    string            `db:"course_id"`
	Description sql.NullString    `db:"description"`
	Status      enum.CommonStatus `db:"status"`
	CreatedBy   string            `db:"created_by"`
	UpdatedBy   sql.NullString    `db:"updated_by"`
	CreatedAt   time.Time         `db:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"`
	DeletedAt   sql.NullTime      `db:"deleted_at"`
}

func (rcv *Syllabus) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		SyllabusFieldID,
		SyllabusFieldTitle,
		SyllabusFieldCode,
		SyllabusFieldTeacherID,
		SyllabusFieldCourseID,
		SyllabusFieldDescription,
		SyllabusFieldStatus,
		SyllabusFieldCreatedBy,
		SyllabusFieldUpdatedBy,
		SyllabusFieldCreatedAt,
		SyllabusFieldUpdatedAt,
		SyllabusFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.Code,
		&rcv.TeacherID,
		&rcv.CourseID,
		&rcv.Description,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Syllabus) TableName() string {
	return "syllabuses"
}
