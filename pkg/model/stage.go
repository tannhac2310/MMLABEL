package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	StageFieldID                = "id"
	StageFieldTitle             = "title"
	StageFieldCourseID          = "course_id"
	StageFieldCategoryID        = "category_id"
	StageFieldStatus            = "status"
	StageFieldTuition           = "tuition"
	StageFieldDiscount          = "discount"
	StageFieldTeacherID         = "teacher_id"
	StageFieldAssistantID       = "assistant_id"
	StageFieldDescription       = "description"
	StageFieldDescriptionShort  = "description_short"
	StageFieldPhoto             = "photo"
	StageFieldVideo             = "video"
	StageFieldCreatedBy         = "created_by"
	StageFieldUpdatedBy         = "updated_by"
	StageFieldCreatedAt         = "created_at"
	StageFieldUpdatedAt         = "updated_at"
	StageFieldDeletedAt         = "deleted_at"
	StageFieldMaxStudent        = "max_student"
	StageFieldCalendar          = "calendar"
	StageFieldProgress          = "progress"
	StageFieldOfficeID          = "office_id"
	StageFieldStageStart        = "stage_start"
	StageFieldStageEnd          = "stage_end"
	StageFieldIsFavorite        = "is_favorite"
	StageFieldDescriptionTarget = "description_target"
	StageFieldCountLesson       = "count_lesson"
	StageFieldSyllabusID        = "syllabus_id"
)

type Calendar struct {
	Day       int
	StartTime string
	EndTime   string
}
type Stage struct {
	ID                string             `db:"id"`
	Title             string             `db:"title"`
	CourseID          string             `db:"course_id"`
	CategoryID        []string           `db:"category_id"`
	Status            enum.StageStatus   `db:"status"`
	Tuition           float64            `db:"tuition"`
	Discount          float64            `db:"discount"`
	TeacherID         []string           `db:"teacher_id"`
	AssistantID       []string           `db:"assistant_id"`
	Description       string             `db:"description"`
	DescriptionShort  string             `db:"description_short"`
	Photo             string             `db:"photo"`
	Video             string             `db:"video"`
	MaxStudent        int                `db:"max_student"`
	Calendar          []*Calendar        `db:"calendar"`
	Progress          float32            `db:"progress"`
	OfficeID          string             `db:"office_id"`
	StageStart        time.Time          `db:"stage_start"`
	StageEnd          time.Time          `db:"stage_end"`
	IsFavorite        enum.CommonBoolean `db:"is_favorite"`
	DescriptionTarget string             `db:"description_target"`
	CountLesson       int                `db:"count_lesson"`
	SyllabusID        string             `db:"syllabus_id"`
	CreatedBy         string             `db:"created_by"`
	UpdatedBy         string             `db:"updated_by"`
	CreatedAt         time.Time          `db:"created_at"`
	UpdatedAt         time.Time          `db:"updated_at"`
	DeletedAt         sql.NullTime       `db:"deleted_at"`
}

func (rcv *Stage) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		StageFieldID,
		StageFieldTitle,
		StageFieldCourseID,
		StageFieldCategoryID,
		StageFieldStatus,
		StageFieldTuition,
		StageFieldDiscount,
		StageFieldTeacherID,
		StageFieldAssistantID,
		StageFieldDescription,
		StageFieldDescriptionShort,
		StageFieldPhoto,
		StageFieldVideo,
		StageFieldCreatedBy,
		StageFieldUpdatedBy,
		StageFieldCreatedAt,
		StageFieldUpdatedAt,
		StageFieldDeletedAt,
		StageFieldMaxStudent,
		StageFieldCalendar,
		StageFieldProgress,
		StageFieldOfficeID,
		StageFieldStageStart,
		StageFieldStageEnd,
		StageFieldIsFavorite,
		StageFieldDescriptionTarget,
		StageFieldCountLesson,
		StageFieldSyllabusID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Title,
		&rcv.CourseID,
		&rcv.CategoryID,
		&rcv.Status,
		&rcv.Tuition,
		&rcv.Discount,
		&rcv.TeacherID,
		&rcv.AssistantID,
		&rcv.Description,
		&rcv.DescriptionShort,
		&rcv.Photo,
		&rcv.Video,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.MaxStudent,
		&rcv.Calendar,
		&rcv.Progress,
		&rcv.OfficeID,
		&rcv.StageStart,
		&rcv.StageEnd,
		&rcv.IsFavorite,
		&rcv.DescriptionTarget,
		&rcv.CountLesson,
		&rcv.SyllabusID,
	}

	return
}

func (*Stage) TableName() string {
	return "stages"
}
