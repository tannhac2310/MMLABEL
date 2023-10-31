package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CommentFieldID        = "id"
	CommentFieldTitle     = "title"
	CommentFieldContent   = "content"
	CommentFieldParentID  = "parent_id"
	CommentFieldProgramID = "program_id"
	CommentFieldStageID   = "stage_id"
	CommentFieldComboID   = "combo_id"
	CommentFieldCourseID  = "course_id"
	CommentFieldLessonID  = "lesson_id"
	CommentFieldStatus    = "status"
	CommentFieldCreatedBy = "created_by"
	CommentFieldUpdatedBy = "updated_by"
	CommentFieldCreatedAt = "created_at"
	CommentFieldUpdatedAt = "updated_at"
	CommentFieldDeletedAt = "deleted_at"
)

type Comment struct {
	ID        string            `db:"id"`
	Title     string            `db:"title"`
	Content   string            `db:"content"`
	ParentID  string            `db:"parent_id"`
	ProgramID string            `db:"program_id"`
	StageID   string            `db:"stage_id"`
	ComboID   string            `db:"combo_id"`
	CourseID  string            `db:"course_id"`
	LessonID  string            `db:"lesson_id"`
	Status    enum.CommonStatus `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
	CreatedBy string            `db:"created_by"`
	UpdatedBy string            `db:"updated_by"`
	UpdatedAt time.Time         `db:"updated_at"`
	DeletedAt sql.NullTime      `db:"deleted_at"`
}

func (c *Comment) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CommentFieldID,
		CommentFieldTitle,
		CommentFieldContent,
		CommentFieldParentID,
		CommentFieldProgramID,
		CommentFieldStageID,
		CommentFieldComboID,
		CommentFieldCourseID,
		CommentFieldLessonID,
		CommentFieldStatus,
		CommentFieldCreatedBy,
		CommentFieldUpdatedBy,
		CommentFieldCreatedAt,
		CommentFieldUpdatedAt,
		CommentFieldDeletedAt,
	}

	values = []interface{}{
		&c.ID,
		&c.Title,
		&c.Content,
		&c.ParentID,
		&c.ProgramID,
		&c.StageID,
		&c.ComboID,
		&c.CourseID,
		&c.LessonID,
		&c.Status,
		&c.CreatedBy,
		&c.UpdatedBy,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	}
	return
}

func (*Comment) TableName() string {
	return "comments"
}
