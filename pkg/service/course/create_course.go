package course

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (c courseService) CreateCourse(ctx context.Context, opt *CreateCourseOpts) (string, error) {
	now := time.Now()

	course := &model.Course{
		ID:                idutil.ULIDNow(),
		Code:              opt.Code,
		Title:             opt.Title,
		CategoryID:        opt.CategoryID,
		Type:              opt.Type,
		Status:            opt.Status,
		Tuition:           opt.Tuition,
		Discount:          opt.Discount,
		TeacherID:         opt.TeacherID,
		AssistantID:       opt.AssistantID,
		Description:       opt.Description,
		Notification:      opt.Notification,
		DescriptionShort:  opt.Description,
		DescriptionTarget: opt.DescriptionTarget,
		Level:             opt.Level,
		IsFavorite:        opt.IsFavorite,
		CountStudent:      opt.CountStudent,
		Photo:             opt.Photo,
		Video:             opt.Video,
		CreatedBy:         opt.CreatedBy,
		UpdatedBy:         opt.CreatedBy,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err := c.courseRepo.Insert(ctx, course)
	if err != nil {
		return "", fmt.Errorf("p.courseRepo.Insert: %w", err)
	}

	return course.ID, nil
}

type CreateCourseOpts struct {
	Code              string
	Title             string
	CategoryID        []string
	Type              enum.CourseType
	Status            enum.CommonStatus
	Tuition           float64
	Discount          float64
	TeacherID         []string
	AssistantID       []string
	Description       string
	DescriptionTarget string
	Notification      *model.Notification
	DescriptionShort  string
	IsFavorite        enum.CommonBoolean
	Level             enum.CourseLevel
	CountStudent      int8
	Photo             string
	Video             string
	CreatedBy         string
}
