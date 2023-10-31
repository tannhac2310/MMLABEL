package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	AttendanceFieldID          = "id"
	AttendanceFieldCourseID    = "course_id"
	AttendanceFieldStageID     = "stage_id"
	AttendanceFieldLessonID    = "lesson_id"
	AttendanceFieldScoreFactor = "score_factor"
	AttendanceFieldNote        = "note"
	AttendanceFieldRecordedAt  = "recorded_at"
	AttendanceFieldCreatedBy   = "created_by"
	AttendanceFieldUpdatedBy   = "updated_by"
	AttendanceFieldCreatedAt   = "created_at"
	AttendanceFieldUpdatedAt   = "updated_at"
	AttendanceFieldDeletedAt   = "deleted_at"
)

type Attendance struct {
	ID          string           `db:"id"`
	CourseID    string           `db:"course_id"`
	StageID     string           `db:"stage_id"`
	LessonID    string           `db:"lesson_id"`
	Note        string           `db:"note"`
	ScoreFactor enum.ScoreFactor `db:"score_factor"`
	RecordedAt  time.Time        `db:"recorded_at"`
	CreatedBy   string           `db:"created_by"`
	UpdatedBy   string           `db:"updated_by"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
	DeletedAt   sql.NullTime     `db:"deleted_at"`
}

func (rcv *Attendance) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		AttendanceFieldID,
		AttendanceFieldCourseID,
		AttendanceFieldStageID,
		AttendanceFieldLessonID,
		AttendanceFieldScoreFactor,
		AttendanceFieldNote,
		AttendanceFieldRecordedAt,
		AttendanceFieldCreatedBy,
		AttendanceFieldUpdatedBy,
		AttendanceFieldCreatedAt,
		AttendanceFieldUpdatedAt,
		AttendanceFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.CourseID,
		&rcv.StageID,
		&rcv.LessonID,
		&rcv.ScoreFactor,
		&rcv.Note,
		&rcv.RecordedAt,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Attendance) TableName() string {
	return "attendances"
}
