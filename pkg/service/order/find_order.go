package order

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/course"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/stage"
)

type DetailInfo struct {
	*model.OrderDetail
	Course *course.Course
	Stage  *stage.Stage
}
type Data struct {
	*repository.OrderData
	DetailInfo []*DetailInfo
}

func (b *orderService) FindOrders(ctx context.Context, opts *FindOrdersOpts, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchOrdersOpts{
		IDs:             opts.IDs,
		StudentID:       opts.StudentID,
		StudentFullName: opts.StudentFullName,
		StudentPhone:    opts.StudentPhone,
		StudentEmail:    opts.StudentEmail,
		InvoiceID:       opts.InvoiceID,
		PaymentMethod:   opts.PaymentMethod,
		SearchString:    opts.SearchString,
		Limit:           limit,
		Offset:          offset,
	}
	orders, err := b.orderRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*Data, 0, len(orders))
	for _, order := range orders {
		orderDetailInfo := make([]*DetailInfo, 0)
		for _, orderDetail := range order.Detail {
			detail := DetailInfo{
				OrderDetail: orderDetail,
			}
			courseData, err := b.courseService.FindCourseByID(ctx, orderDetail.CourseID)
			if err == nil {
				detail.Course = courseData
			}
			stageData, err := b.stageService.FindStageByID(ctx, orderDetail.StageID)
			if err == nil {
				detail.Stage = stageData
			}

			orderDetailInfo = append(orderDetailInfo, &detail)
		}
		data := &Data{
			OrderData:  order,
			DetailInfo: orderDetailInfo,
		}
		results = append(results, data)
	}

	total, err := b.orderRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return results, total, nil
}

type FindOrdersOpts struct {
	IDs             []string
	StudentID       string
	StudentFullName string
	StudentPhone    string
	StudentEmail    string
	InvoiceID       string
	SearchString    string
	PaymentMethod   enum.PaymentMethod
}
