package lesson

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *lessonService) FindLessons(ctx context.Context, opts *FindLessonsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchLessonsOpts{
		IDs:        opts.IDs,
		Search:     opts.Search,
		SyllabusID: opts.SyllabusID,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
	}
	lessons, err := b.lessonRepo.Search(ctx, filter)

	if err != nil {
		return nil, nil, err
	}
	total, err := b.lessonRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*Data, 0, len(lessons))
	for _, lesson := range lessons {
		sInfo, err := b.convertData(ctx, lesson)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, sInfo)
	}
	return results, total, nil
}

type FindLessonsOpts struct {
	IDs        []string
	Search     string
	SyllabusID string
}

func (b *lessonService) convertData(ctx context.Context, lesson *model.Lesson) (*Data, error) {
	sInfo := &Data{
		Lesson: lesson,
	}

	// find syllabus
	syllabus, err := b.syllabusService.FindSyllabusByID(ctx, lesson.SyllabusID)
	if err != nil {
		return nil, fmt.Errorf("lesson.Find.teacher: %w", err)
	}
	sInfo.Syllabus = &model.SyllabusInfo{
		ID:        syllabus.ID,
		Title:     syllabus.Title,
		CourseID:  syllabus.CourseID,
		Course:    syllabus.Course,
		Code:      syllabus.Code,
		TeacherID: syllabus.TeacherID,
		Teacher:   syllabus.Teacher,
		Status:    syllabus.Status,
	}
	// find createdBy
	createdBy, err := b.userRepo.FindByID(ctx, lesson.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("lesson.Find.createdBy: %w", err)
	}
	sInfo.CreatedByInfo = &model.UserInfo{
		Name: createdBy.Name,
		ID:   createdBy.ID,
	}
	// find updatedBy
	if lesson.UpdatedBy.String != "" {
		updatedBy, _ := b.userRepo.FindByID(ctx, lesson.UpdatedBy.String)
		if updatedBy != nil {
			sInfo.UpdatedByInfo = &model.UserInfo{
				Name: updatedBy.Name,
				ID:   updatedBy.ID,
			}
		}
	}
	return sInfo, nil
}
