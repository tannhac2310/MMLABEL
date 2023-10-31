package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	CourseFieldID                = "id"
	CourseFieldCode              = "code"
	CourseFieldTitle             = "title"
	CourseFieldCategoryID        = "category_id"
	CourseFieldType              = "type"
	CourseFieldStatus            = "status"
	CourseFieldTuition           = "tuition"
	CourseFieldDiscount          = "discount"
	CourseFieldTeacherID         = "teacher_id"
	CourseFieldAssistantID       = "assistant_id"
	CourseFieldDescription       = "description"
	CourseFieldDescriptionShort  = "description_short"
	CourseFieldDescriptionTarget = "description_target"
	CourseFieldNotification      = "notification"
	CourseFieldLevel             = "level"
	CourseFieldIsFavorite        = "is_favorite"
	CourseFieldCountStudent      = "count_student"
	CourseFieldPhoto             = "photo"
	CourseFieldVideo             = "video"
	CourseFieldCreatedBy         = "created_by"
	CourseFieldUpdatedBy         = "updated_by"
	CourseFieldCreatedAt         = "created_at"
	CourseFieldUpdatedAt         = "updated_at"
	CourseFieldDeletedAt         = "deleted_at"
)

type Course struct {
	ID                string             `db:"id"`
	Code              string             `db:"code"`
	Title             string             `db:"title"`
	CategoryID        []string           `db:"category_id"`
	Type              enum.CourseType    `db:"type"`
	Status            enum.CommonStatus  `db:"status"`
	Tuition           float64            `db:"tuition"`
	Discount          float64            `db:"discount"`
	TeacherID         []string           `db:"teacher_id"`
	AssistantID       []string           `db:"assistant_id"`
	Description       string             `db:"description"`
	Notification      *Notification      `db:"notification"`
	DescriptionShort  string             `db:"description_short"`
	DescriptionTarget string             `db:"description_target"`
	IsFavorite        enum.CommonBoolean `db:"is_favorite"`
	Level             enum.CourseLevel   `db:"level"`
	CountStudent      int8               `db:"count_student"`
	Photo             string             `db:"photo"`
	Video             string             `db:"video"`
	CreatedBy         string             `db:"created_by"`
	UpdatedBy         string             `db:"updated_by"`
	CreatedAt         time.Time          `db:"created_at"`
	UpdatedAt         time.Time          `db:"updated_at"`
	DeletedAt         sql.NullTime       `db:"deleted_at"`
}

type Notification struct {
	Title        string            `json:"title"`
	Status       enum.CommonStatus `json:"status"`
	Description  string            `json:"description"`
	ButtonTitle  string            `json:"buttonTitle"`
	ButtonLink   string            `json:"buttonLink"`
	ButtonStatus enum.CommonStatus `json:"buttonStatus"`
}

func (c *Course) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CourseFieldID,
		CourseFieldCode,
		CourseFieldTitle,
		CourseFieldCategoryID,
		CourseFieldType,
		CourseFieldStatus,
		CourseFieldLevel,
		CourseFieldTuition,
		CourseFieldDiscount,
		CourseFieldTeacherID,
		CourseFieldAssistantID,
		CourseFieldDescription,
		CourseFieldDescriptionShort,
		CourseFieldDescriptionTarget,
		CourseFieldNotification,
		CourseFieldIsFavorite,
		CourseFieldCountStudent,
		CourseFieldPhoto,
		CourseFieldVideo,
		CourseFieldCreatedBy,
		CourseFieldCreatedAt,
		CourseFieldUpdatedAt,
		CourseFieldUpdatedBy,
		CourseFieldDeletedAt,
	}

	values = []interface{}{
		&c.ID,
		&c.Code,
		&c.Title,
		&c.CategoryID,
		&c.Type,
		&c.Status,
		&c.Level,
		&c.Tuition,
		&c.Discount,
		&c.TeacherID,
		&c.AssistantID,
		&c.Description,
		&c.DescriptionShort,
		&c.DescriptionTarget,
		&c.Notification,
		&c.IsFavorite,
		&c.CountStudent,
		&c.Photo,
		&c.Video,
		&c.CreatedBy,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.UpdatedBy,
		&c.DeletedAt,
	}
	return
}

func (*Course) TableName() string {
	return "courses"
}
