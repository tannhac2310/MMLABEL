package stage

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (s *stageService) EditStage(ctx context.Context, opt *EditStageOpts) error {
	stage, err := s.stageRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("s.pondRepo.FindByID: %w", err)
	}

	stage.Title = opt.Title
	stage.Status = opt.Status
	stage.MaxStudent = opt.MaxStudent
	stage.Tuition = opt.Tuition
	stage.CourseID = opt.CourseID
	stage.SyllabusID = opt.SyllabusID
	stage.CountLesson = opt.CountLesson
	stage.Discount = opt.Discount
	stage.TeacherID = opt.TeacherID
	stage.AssistantID = opt.AssistantID
	stage.Description = opt.Description
	stage.Calendar = opt.Calendar
	stage.DescriptionTarget = opt.DescriptionTarget
	stage.DescriptionShort = opt.DescriptionShort
	stage.IsFavorite = opt.IsFavorite
	stage.Progress = opt.Progress
	stage.Photo = opt.Photo
	stage.Video = opt.Video
	stage.OfficeID = opt.OfficeID
	stage.StageStart = opt.StageStart
	stage.StageEnd = opt.StageEnd
	stage.UpdatedBy = opt.UpdatedBy

	err = s.stageRepo.Update(ctx, stage)
	if err != nil {
		return fmt.Errorf("p.pondRepo.Update: %w", err)
	}

	return nil
}

type EditStageOpts struct {
	ID                string
	Title             string
	CourseID          string
	SyllabusID        string
	CountLesson       int
	Status            enum.StageStatus
	MaxStudent        int
	Tuition           float64
	Discount          float64
	TeacherID         []string
	AssistantID       []string
	Description       string
	Calendar          []*model.Calendar
	DescriptionShort  string
	DescriptionTarget string
	IsFavorite        enum.CommonBoolean
	Progress          float32
	Photo             string
	Video             string
	OfficeID          string
	StageStart        time.Time
	StageEnd          time.Time
	UpdatedBy         string
}
