package lesson

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (b *lessonService) EditLesson(ctx context.Context, opt *EditLessonOpts) error {
	lesson, err := b.lessonRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("b.pondRepo.FindByID: %w", err)
	}

	lesson.Title = opt.Title
	lesson.SyllabusID = opt.SyllabusID
	lesson.Image = cockroach.String(opt.Image)
	lesson.Link = opt.Link
	lesson.Status = opt.Status
	lesson.LessonOrder = opt.LessonOrder
	lesson.Description = cockroach.String(opt.Description)
	lesson.Detail = opt.Detail
	lesson.UpdatedBy = cockroach.String(opt.UpdatedBy)

	err = b.lessonRepo.Update(ctx, lesson)
	if err != nil {
		return fmt.Errorf("p.lessonRepo.Update: %w", err)
	}

	return nil
}

type EditLessonOpts struct {
	ID          string
	Title       string
	SyllabusID  string
	Image       string
	Link        string
	Status      enum.CommonStatus
	LessonOrder int
	Description string
	Detail      []*model.LessonDetail
	UpdatedBy   string
}
