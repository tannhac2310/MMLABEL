package course

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateCourse(ctx context.Context, opt *CreateCourseOpts) (string, error)
	EditCourse(ctx context.Context, opt *EditCourseOpts) error
	FindCourses(ctx context.Context, opts *FindCoursesOpts, limit, offset int64) ([]*Course, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindCourseByID(ctx context.Context, id string) (*Course, error)
}
type Course struct {
	*model.Course
	Teacher    []*model.User
	Assistant  []*model.User
	Categories []*model.Category
}
type courseService struct {
	courseRepo   repository.CourseRepo
	categoryRepo repository.CategoryRepo
	userRepo     repository.UserRepo
}

func NewService(
	courseRepo repository.CourseRepo,
	categoryRepo repository.CategoryRepo,
	userRepo repository.UserRepo,
) Service {
	return &courseService{
		courseRepo:   courseRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}
