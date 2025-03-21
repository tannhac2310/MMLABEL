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
	CountRows(ctx context.Context) (int64, error)
	DeleteByProductionOrderStageID(ctx context.Context, poStageID string) error
	InsertEventLog(ctx context.Context, e *model.EventLog) error
	FindEventLog(ctx context.Context, s *SearchEventLogOpts) ([]*EventLogData, error)
	FindByID(ctx context.Context, id string) (*model.ProductionOrderStageDevice, error)
	GetAssignedByDate(ctx context.Context, dateFrom, dateTo string) ([]model.ProductionOrderStageDevice, error)
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

func (p *productionOrderStageDevicesRepo) GetAssignedByDate(ctx context.Context, dateFrom, dateTo string) ([]model.ProductionOrderStageDevice, error) {
	var result []model.ProductionOrderStageDevice
	sqlQuery := `
		SELECT id, production_order_stage_id, device_id, quantity, settings, estimated_start_at, estimated_complete_at
		FROM production_order_stage_devices 
		WHERE estimated_start_at::DATE = estimated_complete_at::DATE
		AND estimated_start_at::DATE BETWEEN $1 AND $2
		ORDER BY device_id, estimated_start_at
	`
	//if limit > 0 {
	//	sqlQuery += fmt.Sprintf(" LIMIT %d", limit)
	//}
	//if offset > 0 {
	//	sqlQuery += fmt.Sprintf(" OFFSET %d", offset)
	//}
	err := cockroach.Select(ctx, sqlQuery, dateFrom, dateTo).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Query: %w", err)
	}

	return result, nil
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

// CountRows count all rows in the table
func (p *productionOrderStageDevicesRepo) CountRows(ctx context.Context) (int64, error) {
	sqlQuery := `SELECT count(*) as cnt
		FROM production_order_stage_devices`

	countResult := &CountResult{}
	err := cockroach.Select(ctx, sqlQuery).ScanOne(countResult)
	if err != nil {
		return 0, fmt.Errorf("cockroach.Select: %w", err)
	}

	return countResult.Count, nil
}
func (p *productionOrderStageDevicesRepo) DeleteByProductionOrderStageID(ctx context.Context, poStageID string) error {
	sqlQuery := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE production_order_stage_id = $1`

	cmd, err := cockroach.Exec(ctx, sqlQuery, poStageID)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (p *productionOrderStageDevicesRepo) Insert(ctx context.Context, e *model.ProductionOrderStageDevice) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (p *productionOrderStageDevicesRepo) Update(ctx context.Context, e *model.ProductionOrderStageDevice) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (p *productionOrderStageDevicesRepo) SoftDelete(ctx context.Context, id string) error {
	sqlQuery := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sqlQuery, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}
func (p *productionOrderStageDevicesRepo) SoftDeletes(ctx context.Context, ids []string) error {
	sqlQuery := `UPDATE production_order_stage_devices
		SET deleted_at = NOW()
		WHERE id IN ($1)`

	cmd, err := cockroach.Exec(ctx, sqlQuery, strings.Join(ids, ","))
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
	StageIDs                     []string
	Responsible                  []string
	Limit                        int64
	Offset                       int64
	StartAt                      time.Time
	CompleteAt                   time.Time
	EstimatedStartAtFrom         time.Time
	EstimatedStartAtTo           time.Time
	Sort                         *Sort
}

func (s *SearchProductionOrderStageDevicesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ` JOIN devices d ON d.id = b.device_id
		JOIN production_order_stages AS pos ON pos.id = b.production_order_stage_id
		JOIN production_orders AS po ON po.id = pos.production_order_id 
		JOIN stages AS s ON s.id = pos.stage_id 
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

	if len(s.StageIDs) > 0 {
		args = append(args, s.StageIDs)
		conds += fmt.Sprintf(" AND( s.%s = ANY($%d) OR s.parent_id  = ANY($%d))", model.StageFieldID, len(args), len(args))
	}

	if s.StartAt.IsZero() == false {
		args = append(args, s.StartAt)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductionOrderStageDeviceFieldStartAt, len(args))
	}
	if s.CompleteAt.IsZero() == false {
		args = append(args, s.CompleteAt)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductionOrderStageDeviceFieldCompleteAt, len(args))
	}

	if s.EstimatedStartAtFrom.IsZero() == false {
		args = append(args, s.EstimatedStartAtFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductionOrderStageDeviceFieldEstimatedStartAt, len(args))
	}
	if s.EstimatedStartAtTo.IsZero() == false {
		args = append(args, s.EstimatedStartAtTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductionOrderStageDeviceFieldEstimatedStartAt, len(args))
	}

	if len(s.Responsible) > 0 {
		args = append(args, s.Responsible)

		// JOIN production_order_stage_responsible AS posr ON posr.po_stage_device_id = b.id
		joins += ` JOIN production_order_stage_responsible AS posr ON posr.po_stage_device_id = b.id `
		conds += fmt.Sprintf(" AND posr.%s = ANY($%d)", model.ProductionOrderStageResponsibleFieldUserID, len(args))
	}

	b := &model.ProductionOrderStageDevice{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.estimated_start_at DESC "
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

func (p *productionOrderStageDevicesRepo) Search(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) ([]*ProductionOrderStageDeviceData, error) {
	message := make([]*ProductionOrderStageDeviceData, 0)
	sqlQuery, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&message)
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
func (p *productionOrderStageDevicesRepo) FindEventLog(ctx context.Context, s *SearchEventLogOpts) ([]*EventLogData, error) {
	message := make([]*EventLogData, 0)
	sqlQuery, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (p *productionOrderStageDevicesRepo) Count(ctx context.Context, s *SearchProductionOrderStageDevicesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sqlQuery, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sqlQuery, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("productionOrderStageDevicesRepo.Count: %w", err)
	}

	return countResult, nil
}
