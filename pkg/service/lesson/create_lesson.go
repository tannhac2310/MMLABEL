package lesson

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b lessonService) CreateLesson(ctx context.Context, opt *CreateLessonOpts) (string, error) {
	now := time.Now()

	lesson := &model.Lesson{
		ID:          idutil.ULIDNow(),
		Title:       opt.Title,
		SyllabusID:  opt.SyllabusID,
		Image:       cockroach.String(opt.Image),
		Link:        opt.Link,
		Status:      opt.Status,
		LessonOrder: opt.LessonOrder,
		Description: cockroach.String(opt.Description),
		Detail:      opt.Detail,
		CreatedBy:   opt.CreatedBy,
		UpdatedBy:   cockroach.String(opt.CreatedBy),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := b.lessonRepo.Insert(ctx, lesson)
	if err != nil {
		return "", fmt.Errorf("p.lessonRepo.Insert: %w", err)
	}

	return lesson.ID, nil
}

type CreateLessonOpts struct {
	Title       string
	SyllabusID  string
	Image       string
	Link        string
	Status      enum.CommonStatus
	LessonOrder int
	Description string
	Detail      []*model.LessonDetail
	CreatedBy   string
}
