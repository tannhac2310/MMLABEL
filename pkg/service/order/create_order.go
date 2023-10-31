package order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b orderService) CreateOrder(ctx context.Context, opt *CreateOrderOpts) (string, error) {
	now := time.Now()

	order := &model.Order{
		ID:              idutil.ULIDNow(),
		StudentID:       opt.StudentID,
		StudentFullName: opt.StudentFullName,
		StudentPhone:    opt.StudentPhone,
		StudentEmail:    opt.StudentEmail,
		PaymentMethod:   opt.PaymentMethod,
		Address:         opt.Address,
		InvoiceID:       opt.InvoiceID,
		Price:           opt.Price,
		Discount:        opt.Discount,
		Status:          opt.Status,
		Detail:          opt.Detail,
		Note:            opt.Note,
		CreatedAt:       now,
		CreatedBy:       opt.CreatedBy,
		UpdatedBy:       opt.CreatedBy,
		UpdatedAt:       now,
	}
	err := b.orderRepo.Insert(ctx, order)
	if err != nil {
		return "", fmt.Errorf("p.orderRepo.Insert: %w", err)
	}

	return order.ID, nil
}

type CreateOrderOpts struct {
	StudentID       string
	StudentFullName string
	StudentPhone    string
	StudentEmail    string
	PaymentMethod   enum.PaymentMethod
	Address         string
	InvoiceID       string
	Price           float32
	Discount        float32
	DiscountType    enum.DiscountType
	Type            int8
	Status          enum.OrderStatus
	Detail          []*model.OrderDetail
	Note            string
	CreatedBy       string
}
