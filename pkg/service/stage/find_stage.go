package stage

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (s *stageService) FindStages(ctx context.Context, opts *FindStagesOpts, sort *repository.Sort, limit, offset int64) ([]*Stage, *repository.CountResult, error) {
	filter := &repository.SearchStagesOpts{
		IDs:        opts.IDs,
		Title:      opts.Title,
		SyllabusID: opts.SyllabusID,
		CourseID:   opts.CourseID,
		Limit:      limit,
		Offset:     offset,
	}
	stages, err := s.stageRepo.Search(ctx, filter)

	if err != nil {
		return nil, nil, fmt.Errorf("fetch stage: %w", err)
	}
	total, err := s.stageRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*Stage, 0, len(stages))
	for _, stage := range stages {
		// find course
		course, _ := s.courseRepo.FindByID(ctx, stage.CourseID)
		// find teacher
		teachers := make([]*model.User, 0, len(stage.TeacherID))
		for _, teacherID := range stage.TeacherID {
			teacher, err := s.userRepo.FindByID(ctx, teacherID)
			if err == nil {
				teachers = append(teachers, teacher)
			}
		}
		// find assistant
		assistants := make([]*model.User, 0, len(stage.AssistantID))
		for _, assistantID := range stage.TeacherID {
			assistant, err := s.userRepo.FindByID(ctx, assistantID)
			if err == nil {
				assistants = append(assistants, assistant)
			}
		}
		// find syllabus
		syllabus := &model.SyllabusInfo{}
		if stage.SyllabusID != "" {
			rs, err := s.syllabusService.FindSyllabusByID(ctx, stage.SyllabusID)
			if err == nil {
				syllabus = &model.SyllabusInfo{
					ID:        rs.ID,
					Title:     rs.Title,
					CourseID:  rs.CourseID,
					Course:    rs.Course,
					Code:      rs.Code,
					TeacherID: rs.TeacherID,
					Teacher:   rs.Teacher,
					Status:    rs.Status,
				}
			}
		}
		for _, assistantID := range stage.TeacherID {
			assistant, err := s.userRepo.FindByID(ctx, assistantID)
			if err == nil {
				assistants = append(assistants, assistant)
			}
		}
		// count student in class
		count, err := s.stageStudentRepo.Count(ctx, &repository.SearchStageStudentsOpts{
			CourseID: stage.CourseID,
			StageID:  stage.ID,
		})
		if err != nil {
			return nil, nil, err
		}

		results = append(results, &Stage{
			CountStageStudent: count.Count,
			Stage:             stage,
			Course:            course,
			Teacher:           teachers,
			Assistant:         assistants,
			Syllabus:          syllabus,
		})
	}
	return results, total, nil
}

type FindStagesOpts struct {
	IDs        []string
	Title      string
	CourseID   string
	SyllabusID string
}
