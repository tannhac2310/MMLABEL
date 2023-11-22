package repository

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	model2 "mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type ProductionOrderStageDeviceRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrderStageDevice) error
	Update(ctx context.Context, e *model.ProductionOrderStageDevice) error
	SoftDelete(ctx context.Context, id string) error
	SoftDeletes(ctx context.Context, ids []string) error
	Search(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) ([]*ProductionOrderStageDeviceData, error)
	Count(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) (*CountResult, error)
	DeleteByProductionOrderStageID(ctx context.Context, poStageID string) error
	InsertEventLog(ctx context.Context, e *model.EventLog) error
	FindEventLog(ctx context.Context, s *SearchEventLogOpts) ([]*EventLogData, error)
}
type SearchEventLogOpts struct {
	DeviceID string
	Date     string
}
type EventLogData struct {
	*model.EventLog
	DeviceName string `db:"device_name"`
}
type ProductionOrderStageDeviceData struct {
	*model.ProductionOrderStageDevice
	DeviceName        string `db:"device_name"`
	ProductionOrderID string `db:"production_order_id"`
	ResponsibleObject []*model2.User
}

type productionOrderStageDevicesRepo struct {
}

func NewProductionOrderStageDeviceRepo() ProductionOrderStageDeviceRepo {
	return &productionOrderStageDevicesRepo{}
}

func (p *productionOrderStageDevicesRepo) InsertEventLog(ctx context.Context, e *model.EventLog) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}
func (p *productionOrderStageDevicesRepo) DeleteByProductionOrderStageID(ctx context.Context, poStageID string) error {
	sql := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE production_order_stage_id = $1`

	cmd, err := cockroach.Exec(ctx, sql, poStageID)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *productionOrderStageDevicesRepo) Insert(ctx context.Context, e *model.ProductionOrderStageDevice) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *productionOrderStageDevicesRepo) Update(ctx context.Context, e *model.ProductionOrderStageDevice) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *productionOrderStageDevicesRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}
func (r *productionOrderStageDevicesRepo) SoftDeletes(ctx context.Context, ids []string) error {
	sql := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE id IN ($1)`

	cmd, err := cockroach.Exec(ctx, sql, strings.Join(ids, ","))
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}
	return nil
}

// SearchProductionOrderStageDevicesOpts all params is options
type SearchProductionOrderStageDevicesOpts struct {
	ProductionOrderStageID     string
	ProductionOrderID          string
	DeviceID                   string
	ProductionOrderStageStatus enum.ProductionOrderStageStatus
	Status                     enum.ProductionOrderStageDeviceStatus
	Limit                      int64
	Offset                     int64
	Sort                       *Sort
}

func (s *SearchProductionOrderStageDevicesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND pos.production_order_id = $%d", len(args))
	}

	if s.ProductionOrderStageID != "" {
		args = append(args, s.ProductionOrderStageID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageDeviceFieldProductionOrderStageID, len(args))
	}

	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageDeviceFieldDeviceID, len(args))
	}

	if s.Status > 0 {
		args = append(args, s.Status)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageDeviceFieldStatus, len(args))
	}
	if s.ProductionOrderStageStatus > 0 {
		args = append(args, s.ProductionOrderStageStatus)
		conds += fmt.Sprintf(" AND pos.%s = $%d", model.ProductionOrderStageFieldStatus, len(args))
	}

	b := &model.ProductionOrderStageDevice{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s, pos.production_order_id as production_order_id, COALESCE (d.name,'N/A') as device_name
		FROM %s AS b %s
		JOIN devices d ON d.id = b.device_id
		JOIN production_order_stages AS pos ON pos.id = b.production_order_stage_id
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

func (r *productionOrderStageDevicesRepo) Search(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) ([]*ProductionOrderStageDeviceData, error) {
	message := make([]*ProductionOrderStageDeviceData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}
func (s *SearchEventLogOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND el.device_id = $%d", len(args))
	}
	if s.Date != "" {
		args = append(args, s.Date)
		conds += fmt.Sprintf(" AND el.date = $%d", len(args))
	}

	b := &model.EventLog{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS el %s
		WHERE TRUE %s AND el.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY el.id DESC "
	return fmt.Sprintf(`SELECT el.%s, COALESCE (d.name,'N/A') as device_name
		FROM %s AS el %s
		LEFT JOIN devices d ON d.id = el.device_id
		WHERE TRUE %s
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", el."), b.TableName(), joins, conds, order, 1000, 0), args

}
func (r *productionOrderStageDevicesRepo) FindEventLog(ctx context.Context, s *SearchEventLogOpts) ([]*EventLogData, error) {
	message := make([]*EventLogData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *productionOrderStageDevicesRepo) Count(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
