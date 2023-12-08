package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DeviceWorkingHistoryRepo interface {
	Insert(ctx context.Context, e *model.DeviceWorkingHistory) error
	Update(ctx context.Context, e *model.DeviceWorkingHistory) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.DeviceWorkingHistory, error)
	Search(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) ([]*DeviceWorkingHistoryData, error)
	Count(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) (*CountResult, error)
}

type sDeviceWorkingHistoryRepo struct {
}

func NewDeviceWorkingHistoryRepo() DeviceWorkingHistoryRepo {
	return &sDeviceWorkingHistoryRepo{}
}

func (r *sDeviceWorkingHistoryRepo) Insert(ctx context.Context, e *model.DeviceWorkingHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sDeviceWorkingHistoryRepo) FindByID(ctx context.Context, id string) (*model.DeviceWorkingHistory, error) {
	e := &model.DeviceWorkingHistory{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sDeviceWorkingHistoryRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sDeviceWorkingHistoryRepo) Update(ctx context.Context, e *model.DeviceWorkingHistory) error {
	e.UpdatedAt = cockroach.Time(time.Now())
	return cockroach.Update(ctx, e)
}

func (r *sDeviceWorkingHistoryRepo) SoftDelete(ctx context.Context, id string) error {
	return fmt.Errorf("not implement")
}

// SearchDeviceWorkingHistoryOpts all params is options
type SearchDeviceWorkingHistoryOpts struct {
	IDs      []string
	DeviceID string
	Date     string
	// todo add more search options
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
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DeviceWorkingHistoryFieldID)
	}

	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceWorkingHistoryFieldDeviceID, len(args))
	}

	if s.Date != "" {
		args = append(args, s.Date)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceWorkingHistoryFieldDate, len(args))
	}

	b := &model.DeviceWorkingHistory{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s  ", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s  %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type DeviceWorkingHistoryData struct {
	*model.DeviceWorkingHistory
}

func (r *sDeviceWorkingHistoryRepo) Search(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) ([]*DeviceWorkingHistoryData, error) {
	DeviceWorkingHistory := make([]*DeviceWorkingHistoryData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&DeviceWorkingHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return DeviceWorkingHistory, nil
}

func (r *sDeviceWorkingHistoryRepo) Count(ctx context.Context, s *SearchDeviceWorkingHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sDeviceWorkingHistoryRepo.Count: %w", err)
	}

	return countResult, nil
}
