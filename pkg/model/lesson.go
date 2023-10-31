package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	LessonFieldID          = "id"
	LessonFieldTitle       = "title"
	LessonFieldSyllabusID  = "syllabus_id"
	LessonFieldImage       = "image"
	LessonFieldLink        = "link"
	LessonFieldStatus      = "status"
	LessonFieldLessonOrder = "lesson_order"
	LessonFieldDescription = "description"
	LessonFieldDetail      = "detail"
	LessonFieldCreatedBy   = "created_by"
	LessonFieldUpdatedBy   = "updated_by"
	LessonFieldCreatedAt   = "created_at"
	LessonFieldUpdatedAt   = "updated_at"
	LessonFieldDeletedAt   = "deleted_at"
)

type LessonDetail struct {
	Title       string  `db:"title"`
	TimeSpent   float64 `db:"timeSpent"`
	DetailOrder int8    `db:"detailOrder"`
	Description string  `db:"description"`
}
type Lesson struct {
	ID          string            `db:"id"`
	Title       string            `db:"title"`
	SyllabusID  string            `db:"syllabus_id"`
	Image       sql.NullString    `db:"image"`
	Link        string            `db:"link"`
	Status      enum.CommonStatus `db:"status"`
	LessonOrder int               `db:"lesson_order"`
	Description sql.NullString    `db:"description"`
	Detail      []*LessonDetail   `db:"detail"`
	CreatedBy   string            `db:"created_by"`
	UpdatedBy   sql.NullString    `db:"updated_by"`
	CreatedAt   time.Time         `db:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"`
	DeletedAt   sql.NullTime      `db:"deleted_at"`
}

func (rcv *Lesson) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		LessonFieldID,
		LessonFieldTitle,
		LessonFieldSyllabusID,
		LessonFieldImage,
		LessonFieldLink,
		LessonFieldStatus,
		LessonFieldLessonOrder,
		LessonFieldDescription,
		LessonFieldDetail,
		LessonFieldCreatedBy,
		LessonFieldUpdatedBy,
		LessonFieldCreatedAt,
		LessonFieldUpdatedAt,
		LessonFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.SyllabusID,
		&rcv.Image,
		&rcv.Link,
		&rcv.Status,
		&rcv.LessonOrder,
		&rcv.Description,
		&rcv.Detail,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Lesson) TableName() string {
	return "lessons"
}
