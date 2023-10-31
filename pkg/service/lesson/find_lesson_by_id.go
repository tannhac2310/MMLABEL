package lesson

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *lessonService) FindLessonByID(ctx context.Context, id string) (*Data, error) {
	lesson, err := b.lessonRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sInfo := &Data{
		Lesson: lesson,
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
