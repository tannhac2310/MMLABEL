package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DeviceProgressStatusHistoryRepo interface {
	Insert(ctx context.Context, e *model.DeviceProgressStatusHistory) error
	Update(ctx context.Context, e *model.DeviceProgressStatusHistory) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) ([]*DeviceProgressStatusHistoryData, error)
	Count(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) (*CountResult, error)
	FindProductionOrderStageDeviceID(ctx context.Context, ProductionOrderStageID string, deviceID string) (*DeviceProgressStatusHistoryData, error)
}

type sDeviceProgressStatusHistoryRepo struct {
}

func NewDeviceProgressStatusHistoryRepo() DeviceProgressStatusHistoryRepo {
	return &sDeviceProgressStatusHistoryRepo{}
}
func (i *sDeviceProgressStatusHistoryRepo) FindProductionOrderStageDeviceID(ctx context.Context, ProductionOrderStageID string, deviceID string) (*DeviceProgressStatusHistoryData, error) {
	deviceProcessStatusHistoryData := &DeviceProgressStatusHistoryData{}
	sql := `SELECT * FROM device_progress_status_history WHERE production_order_stage_device_id = $1 AND device_id = $2 AND is_resolved = 0`
	err := cockroach.Select(ctx, sql, ProductionOrderStageID, deviceID).ScanOne(deviceProcessStatusHistoryData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return deviceProcessStatusHistoryData, nil
}
func (r *sDeviceProgressStatusHistoryRepo) Insert(ctx context.Context, e *model.DeviceProgressStatusHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sDeviceProgressStatusHistoryRepo) Update(ctx context.Context, e *model.DeviceProgressStatusHistory) error {
	return cockroach.Update(ctx, e)
}

func (r *sDeviceProgressStatusHistoryRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE device_progress_status_history SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("device_progress_status_history cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sDeviceProgressStatusHistoryRepo not found any records to delete")
	}

	return nil
}

// SearchDeviceProgressStatusHistoryOpts all params is options
type SearchDeviceProgressStatusHistoryOpts struct {
	IDs         []string
	CreatedFrom time.Time
	CreatedTo   time.Time
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchDeviceProgressStatusHistoryOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DeviceProgressStatusHistoryFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.DeviceProgressStatusHistoryFieldName, model.DeviceProgressStatusHistoryFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.DeviceProgressStatusHistoryFieldCode, len(args))
	//}
	if s.CreatedFrom.IsZero() == false {
		args = append(args, s.CreatedFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.DeviceProgressStatusHistoryFieldCreatedAt, len(args))
	}
	if s.CreatedTo.IsZero() == false {
		args = append(args, s.CreatedTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.DeviceProgressStatusHistoryFieldCreatedAt, len(args))
	}

	b := &model.DeviceProgressStatusHistory{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type DeviceProgressStatusHistoryData struct {
	*model.DeviceProgressStatusHistory
}

func (r *sDeviceProgressStatusHistoryRepo) Search(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) ([]*DeviceProgressStatusHistoryData, error) {
	DeviceProgressStatusHistory := make([]*DeviceProgressStatusHistoryData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&DeviceProgressStatusHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return DeviceProgressStatusHistory, nil
}

func (r *sDeviceProgressStatusHistoryRepo) Count(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sDeviceProgressStatusHistoryRepo.Count: %w", err)
	}

	return countResult, nil
}
