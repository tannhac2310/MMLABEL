package lesson

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/syllabus"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateLesson(ctx context.Context, opt *CreateLessonOpts) (string, error)
	EditLesson(ctx context.Context, opt *EditLessonOpts) error
	FindLessons(ctx context.Context, opts *FindLessonsOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindLessonByID(ctx context.Context, id string) (*Data, error)
}

type lessonService struct {
	lessonRepo      repository.LessonRepo
	userRepo        repository.UserRepo
	syllabusService syllabus.Service
}

func NewService(
	lessonRepo repository.LessonRepo,
	userRepo repository.UserRepo,
	syllabusService syllabus.Service,
) Service {
	return &lessonService{
		lessonRepo:      lessonRepo,
		userRepo:        userRepo,
		syllabusService: syllabusService,
	}
}

type Data struct {
	*model.Lesson
	Syllabus      *model.SyllabusInfo
	CreatedByInfo *model.UserInfo
	UpdatedByInfo *model.UserInfo
}
