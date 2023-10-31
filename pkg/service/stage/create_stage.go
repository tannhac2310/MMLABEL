package stage

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (s stageService) CreateStage(ctx context.Context, opt *CreateStageOpts) (string, error) {
	now := time.Now()
	stage := &model.Stage{
		ID:                idutil.ULIDNow(),
		Title:             opt.Title,
		Status:            opt.Status,
		CourseID:          opt.CourseID,
		SyllabusID:        opt.SyllabusID,
		MaxStudent:        opt.MaxStudent,
		Tuition:           opt.Tuition,
		Discount:          opt.Discount,
		TeacherID:         opt.TeacherID,
		AssistantID:       opt.AssistantID,
		Calendar:          opt.Calendar,
		Description:       opt.Description,
		CountLesson:       opt.CountLesson,
		DescriptionShort:  opt.DescriptionShort,
		DescriptionTarget: opt.DescriptionTarget,
		IsFavorite:        opt.IsFavorite,
		Progress:          opt.Progress,
		Photo:             opt.Photo,
		Video:             opt.Video,
		OfficeID:          opt.OfficeID,
		StageStart:        opt.StageStart,
		StageEnd:          opt.StageEnd,
		CreatedBy:         opt.CreatedBy,
		UpdatedBy:         opt.CreatedBy,
		UpdatedAt:         now,
		CreatedAt:         now,
	}
	err := s.stageRepo.Insert(ctx, stage)
	if err != nil {
		return "", fmt.Errorf("p.stageRepo.Insert: %w", err)
	}

	return stage.ID, nil
}

type CreateStageOpts struct {
	Title             string
	Status            enum.StageStatus
	MaxStudent        int
	CountLesson       int
	CourseID          string
	SyllabusID        string
	Tuition           float64
	Discount          float64
	TeacherID         []string
	AssistantID       []string
	Description       string
	DescriptionTarget string
	Calendar          []*model.Calendar
	DescriptionShort  string
	IsFavorite        enum.CommonBoolean
	Progress          float32
	Photo             string
	Video             string
	OfficeID          string
	StageStart        time.Time
	StageEnd          time.Time
	CreatedBy         string
}
