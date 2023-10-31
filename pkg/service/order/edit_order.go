package order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b *orderService) EditOrder(ctx context.Context, opt *EditOrderOpts) error {
	table := model.Order{}
	updater := cockroach.NewUpdater(table.TableName(), model.OrderFieldID, opt.ID)

	updater.Set(model.OrderFieldStudentID, opt.StudentID)
	updater.Set(model.OrderFieldStudentFullName, opt.StudentFullName)
	updater.Set(model.OrderFieldStudentPhone, opt.StudentPhone)
	updater.Set(model.OrderFieldStudentEmail, opt.StudentEmail)
	updater.Set(model.OrderFieldPaymentMethod, opt.PaymentMethod)
	updater.Set(model.OrderFieldAddress, opt.Address)
	updater.Set(model.OrderFieldPrice, opt.Price)
	updater.Set(model.OrderFieldDiscount, opt.Discount)
	updater.Set(model.OrderFieldStatus, opt.Status)
	updater.Set(model.OrderFieldDetail, opt.Detail)
	updater.Set(model.OrderFieldNote, opt.Note)

	updater.Set(model.OrderFieldUpdatedAt, time.Now())
	updater.Set(model.OrderFieldUpdatedBy, opt.UpdatedBy)

	err := cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update order failed %w", err)
	}
	return nil
}

type EditOrderOpts struct {
	ID              string
	StudentID       string
	StudentFullName string
	StudentPhone    string
	StudentEmail    string
	PaymentMethod   enum.PaymentMethod
	Address         string
	Price           float32
	Discount        float32
	Status          enum.OrderStatus
	Detail          []*model.OrderDetail
	Note            string
	UpdatedBy       string
}
