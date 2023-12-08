package repository

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/type/date"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DeviceWorkingHistoryRepo interface {
	Insert(ctx context.Context, e *model.DeviceWorkingHistory) error
	Update(ctx context.Context, e *model.DeviceWorkingHistory) error
	FindByDeviceAndDate(ctx context.Context, device_id string, working_date date) error
	Search(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) ([]*DeviceData, error)
	Count(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) (*CountResult, error)
}

type deviceWorkingHistoryRepo struct {
}

func NewDeviceWorkingHistoryRepo() DeviceWorkingHistoryRepo {
	return &deviceWorkingHistoryRepo{}
}

func (r *deviceWorkingHistoryRepo) Insert(ctx context.Context, e *model.DeviceWorkingHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *deviceWorkingHistoryRepo) Update(ctx context.Context, e *model.DeviceWorkingHistory) error {
	
	return cockroach.Update(ctx, e)
}

func (r *deviceWorkingHistoryRepo) FindByDeviceAndDate(ctx context.Context, device_id string, working_date date.Date) error {

}


// SearchDeviceWorkingHistoryOpts all params is options
type SearchDeviceWorkingHistoryOpts struct {
	IDs    						   []string
	ProductionOrderStageDeviceID 	string
	DeviceID 					    string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchDeviceWorkingHistoryOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DeviceFieldID)
	}
	
	if len(s.ProductionOrderStageDeviceID) > 0 {
		args = append(args, s.ProductionOrderStageDeviceID)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceWorkingHistoryFieldProductionOrderStageDeviceID, len(args))
	}
	
	if len(s.DeviceID) > 0 {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceWorkingHistoryFieldDeviceID, len(args))
	}
	

	b := &model.DeviceWorkingHistory{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s `, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s 
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type DeviceWorkingHistoryData struct {
	*model.DeviceWorkingHistory
}

func (r *deviceWorkingHistoryRepo) Search(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) ([]*DeviceWorkingHistoryData, error) {
	devicesWorkingHistory := make([]*DeviceWorkingHistoryData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&devicesWorkingHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return devicesWorkingHistory, nil
}

func (r *deviceWorkingHistoryRepo) Count(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
