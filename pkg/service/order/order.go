package order

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/stage"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/course"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateOrder(ctx context.Context, opt *CreateOrderOpts) (string, error)
	EditOrder(ctx context.Context, opt *EditOrderOpts) error
	FindOrders(ctx context.Context, opts *FindOrdersOpts, limit, offset int64) ([]*Data, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindOrderByID(ctx context.Context, id string) (*Data, error)
}

type orderService struct {
	orderRepo     repository.OrderRepo
	courseService course.Service
	stageService  stage.Service
}

func NewService(
	orderRepo repository.OrderRepo,
	courseService course.Service,
	stageService stage.Service,
) Service {
	return &orderService{
		orderRepo:     orderRepo,
		courseService: courseService,
		stageService:  stageService,
	}
}
