package course

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (c *courseService) EditCourse(ctx context.Context, opt *EditCourseOpts) error {
	course, err := c.courseRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("EditCourse:c.courseRepo.FindByID %w, %s", err, opt.ID)
	}

	course.Code = opt.Code
	course.Title = opt.Title
	course.CategoryID = opt.CategoryID
	course.Type = opt.Type
	course.Status = opt.Status
	course.Level = opt.Level
	course.Tuition = opt.Tuition
	course.Discount = opt.Discount
	course.TeacherID = opt.TeacherID
	course.AssistantID = opt.AssistantID
	course.Description = opt.Description
	course.Notification = opt.Notification
	course.DescriptionShort = opt.DescriptionShort
	course.IsFavorite = opt.IsFavorite
	course.CountStudent = opt.CountStudent
	course.Photo = opt.Photo
	course.DescriptionTarget = opt.DescriptionTarget

	err = c.courseRepo.Update(ctx, course)
	if err != nil {
		return fmt.Errorf("EditCourse:p.pondRepo.Update: %w", err)
	}

	return nil
}

type EditCourseOpts struct {
	ID                string
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
}
