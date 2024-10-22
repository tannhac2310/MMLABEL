package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

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
	FindByID(ctx context.Context, id string) (*model.ProductionOrderStageDevice, error)
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
	DeviceName                              string                          `db:"device_name"`
	DeviceData                              any                             `db:"device_data"`
	ProductionOrderID                       string                          `db:"production_order_id"`
	ProductionOrderName                     string                          `db:"production_order_name"`
	ProductionOrderStatus                   enum.ProductionOrderStatus      `db:"production_order_status"`
	ProductionOrderStageName                string                          `db:"production_order_stage_name"`
	ProductionOrderStageCode                string                          `db:"production_order_stage_code"`
	ProductionOrderStageStatus              enum.ProductionOrderStageStatus `db:"production_order_stage_status"`
	ProductionOrderStageStartedAt           sql.NullTime                    `db:"production_order_stage_started_at"`
	ProductionOrderStageCompletedAt         sql.NullTime                    `db:"production_order_stage_completed_at"`
	ProductionOrderStageEstimatedStartAt    sql.NullTime                    `db:"production_order_stage_estimated_start_at"`
	ProductionOrderStageEstimatedCompleteAt sql.NullTime                    `db:"production_order_stage_estimated_complete_at"`
	ResponsibleObject                       []*model2.User
}

type productionOrderStageDevicesRepo struct {
}

func NewProductionOrderStageDeviceRepo() ProductionOrderStageDeviceRepo {
	return &productionOrderStageDevicesRepo{}
}

func (p *productionOrderStageDevicesRepo) FindByID(ctx context.Context, id string) (*model.ProductionOrderStageDevice, error) {
	e := &model.ProductionOrderStageDevice{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
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
	ID                           string
	IDs                          []string
	ProductionOrderStageIDs      []string
	ProductionOrderIDs           []string
	ProcessStatuses              []enum.ProductionOrderStageDeviceStatus
	DeviceIDs                    []string
	ProductionOrderStageStatuses []enum.ProductionOrderStageStatus
	Responsible                  []string
	Limit                        int64
	Offset                       int64
	Sort                         *Sort
}

func (s *SearchProductionOrderStageDevicesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ` JOIN devices d ON d.id = b.device_id
		JOIN production_order_stages AS pos ON pos.id = b.production_order_stage_id
		JOIN production_orders AS po ON po.id = pos.production_order_id 
		JOIN stages AS s ON s.id = pos.stage_id 
		JOIN production_order_stage_responsible AS posr ON posr.po_stage_device_id = b.id
`

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderStageDeviceFieldID, len(args))
	}
	if s.ID != "" {
		args = append(args, s.ID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageDeviceFieldID, len(args))
	}
	if len(s.ProductionOrderStageIDs) > 0 {
		args = append(args, s.ProductionOrderStageIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderStageDeviceFieldProductionOrderStageID, len(args))
	}
	if len(s.ProductionOrderIDs) > 0 {
		args = append(args, s.ProductionOrderIDs)
		conds += fmt.Sprintf(" AND pos.%s = ANY($%d)", model.ProductionOrderStageFieldProductionOrderID, len(args))
	}
	if len(s.ProcessStatuses) > 0 {
		args = append(args, s.ProcessStatuses)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderStageDeviceFieldProcessStatus, len(args))
	}
	if len(s.DeviceIDs) > 0 {
		args = append(args, s.DeviceIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.ProductionOrderStageDeviceFieldDeviceID, len(args))
	}
	if len(s.ProductionOrderStageStatuses) > 0 {
		args = append(args, s.ProductionOrderStageStatuses)
		conds += fmt.Sprintf(" AND pos.%s = ANY($%d)", model.ProductionOrderStageFieldStatus, len(args))
	}
	if len(s.Responsible) > 0 {
		args = append(args, s.Responsible)
		conds += fmt.Sprintf(" AND posr.%s = ANY($%d)", model.ProductionOrderStageResponsibleFieldUserID, len(args))
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
		order = fmt.Sprintf(" ORDER BY b.%s %s ", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s, pos.production_order_id as production_order_id, COALESCE (d.name,'N/A') as device_name, d.data as device_data,
		po.name as production_order_name, po.status as production_order_status, 
    	pos.status as production_order_stage_status,
    	pos.started_at as production_order_stage_started_at, pos.completed_at as production_order_stage_completed_at,
    	pos.estimated_start_at as production_order_stage_estimated_start_at,
    	pos.estimated_complete_at as production_order_stage_estimated_complete_at,
    	s.code as production_order_stage_code,
		s.name as production_order_stage_name
		FROM %s AS b %s
		
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
		OFFSET %d`, strings.Join(fields, ", el."), b.TableName(), joins, conds, order, 50000, 0), args

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
		return nil, fmt.Errorf("productionOrderStageDevicesRepo.Count: %w", err)
	}

	return countResult, nil
}
