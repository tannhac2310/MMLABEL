package product

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateProduct(ctx context.Context, opt *CreateProductOpts) (string, error)
	UpdateProduct(ctx context.Context, opt *UpdateProductOpts) error
	DeleteProduct(ctx context.Context, id string) error
	FindProduct(ctx context.Context, opts *FindProductOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
}

type Data struct {
	ID                 string
	Name               string
	Code               string
	CustomerID         string
	CustomerData       *repository.CustomerData
	ProductionPlanData *repository.ProductionPlanData
	ProductionPlanID   string
	SaleID             string
	Description        string
	Data               any
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreatedBy          string
	CreatedByName      string
	UpdatedBy          string
	UpdatedByName      string
}
type CreateProductOpts struct {
	Name        string
	Code        string
	CustomerID  string
	SaleID      string
	Description string
	Data        any
	CreatedBy   string
	UserField   []*UserField
}

type UserField struct {
	Key   string
	Value string
}

type UpdateProductOpts struct {
	ID          string
	Name        string
	Code        string
	CustomerID  string
	SaleID      string
	Description string
	Data        any
	UpdatedBy   string
	UserField   []*UserField
}

type FindProductOpts struct {
	Name           string
	Code           string
	CustomerID     string
	SaleID         string
	ProductPlanID  string
	ProductOrderID string
}
type productService struct {
	productRepo        repository.ProductRepo
	userFieldRepo      repository.CustomFieldRepo
	customerRepo       repository.CustomerRepo
	userRepo           repository2.UserRepo
	productionPlanRepo repository.ProductionPlanRepo
}

func NewService(productRepo repository.ProductRepo, userFieldRepo repository.CustomFieldRepo, customerRepo repository.CustomerRepo, userRepo repository2.UserRepo, productionPlanRepo repository.ProductionPlanRepo) Service {
	return &productService{
		productRepo:        productRepo,
		userFieldRepo:      userFieldRepo,
		customerRepo:       customerRepo,
		userRepo:           userRepo,
		productionPlanRepo: productionPlanRepo,
	}
}
