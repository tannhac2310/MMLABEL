package syllabus

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *syllabusService) FindSyllabuses(ctx context.Context, opts *FindSyllabusesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchSyllabusesOpts{
		IDs:       opts.IDs,
		Search:    opts.Search,
		CourseID:  opts.CourseID,
		TeacherID: opts.TeacherID,
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
	}
	syllabuses, err := b.syllabusRepo.Search(ctx, filter)

	if err != nil {
		return nil, nil, err
	}
	total, err := b.syllabusRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*Data, 0, len(syllabuses))
	for _, syllabus := range syllabuses {
		sInfo, err := b.convertData(ctx, syllabus)
		if err != nil {
			return nil, nil, err
		}
		results = append(results, sInfo)
	}
	return results, total, nil
}

type FindSyllabusesOpts struct {
	IDs       []string
	Search    string
	TeacherID string
	CourseID  string
}

func (b *syllabusService) convertData(ctx context.Context, syllabus *model.Syllabus) (*Data, error) {
	sInfo := &Data{
		Syllabus: syllabus,
	}
	// find teacher
	teacher, err := b.userRepo.FindByID(ctx, syllabus.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("syllabus.Find.teacher: %w", err)
	}
	sInfo.Teacher = &model.UserInfo{
		Name: teacher.Name,
		ID:   teacher.ID,
	}
	// find course
	course, err := b.courseRepo.FindByID(ctx, syllabus.CourseID)
	if err != nil {
		return nil, fmt.Errorf("syllabus.Find.teacher: %w", err)
	}
	sInfo.Course = &model.CourseInfo{
		ID:     course.ID,
		Code:   course.Code,
		Title:  course.Title,
		Type:   course.Type,
		Status: course.Status,
	}
	// find createdBy
	createdBy, err := b.userRepo.FindByID(ctx, syllabus.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("syllabus.Find.createdBy: %w", err)
	}
	sInfo.CreatedByInfo = &model.UserInfo{
		Name: createdBy.Name,
		ID:   createdBy.ID,
	}
	// find updatedBy
	if syllabus.UpdatedBy.String != "" {
		updatedBy, _ := b.userRepo.FindByID(ctx, syllabus.UpdatedBy.String)
		if updatedBy != nil {
			sInfo.UpdatedByInfo = &model.UserInfo{
				Name: updatedBy.Name,
				ID:   updatedBy.ID,
			}
		}
	}
	return sInfo, nil
}
