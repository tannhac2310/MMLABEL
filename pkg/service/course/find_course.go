package course

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (c *courseService) FindCourses(ctx context.Context, opts *FindCoursesOpts, limit, offset int64) ([]*Course, *repository.CountResult, error) {
	filter := &repository.SearchCoursesOpts{
		IDs:    opts.IDs,
		Code:   opts.Code,
		Search: opts.Search,
		Title:  opts.Title,
		Type:   opts.Type,
		Status: opts.Status,
		Limit:  limit,
		Offset: offset,
	}
	courses, err := c.courseRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := c.courseRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*Course, 0, len(courses))
	for _, course := range courses {
		cats := make([]*model.Category, 0, len(course.CategoryID))
		for _, catID := range course.CategoryID {
			category, err := c.categoryRepo.FindByID(ctx, catID)
			if err == nil {
				cats = append(cats, category)
			}
		}
		// find teacher
		teachers := make([]*model.User, 0, len(course.TeacherID))
		for _, teacherID := range course.TeacherID {
			teacher, err := c.userRepo.FindByID(ctx, teacherID)
			if err == nil {
				teachers = append(teachers, teacher)
			}
		}
		// find assistant
		assistants := make([]*model.User, 0, len(course.AssistantID))
		for _, assistantID := range course.TeacherID {
			assistant, err := c.userRepo.FindByID(ctx, assistantID)
			if err == nil {
				assistants = append(assistants, assistant)
			}
		}
		results = append(results, &Course{
			Course:     course,
			Teacher:    teachers,
			Assistant:  assistants,
			Categories: cats,
		})
	}
	return results, total, nil
}

type FindCoursesOpts struct {
	IDs    []string
	Code   string
	Search string
	Title  string
	Type   enum.CourseType
	Status enum.CommonStatus
}
