package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type DeviceBrokenHistoryRepo interface {
	Insert(ctx context.Context, e *model.DeviceBrokenHistory) error
	Update(ctx context.Context, e *model.DeviceBrokenHistory) error
	Search(ctx context.Context, s *SearchDeviceBrokenHistoryOpts) ([]*DeviceBrokenHistoryData, error)
	Count(ctx context.Context, s *SearchDeviceBrokenHistoryOpts) (*CountResult, error)
}

type sDeviceBrokenHistoryRepo struct {
}

func NewDeviceBrokenHistoryRepo() DeviceBrokenHistoryRepo {
	return &sDeviceBrokenHistoryRepo{}
}

func (i *sDeviceBrokenHistoryRepo) FindByID(ctx context.Context, ID string) (*DeviceBrokenHistoryData, error) {
	deviceProcessStatusHistoryData := &DeviceBrokenHistoryData{}
	sqlQuery := `SELECT * FROM device_progress_status_history WHERE ID = $1 ORDER BY ID DESC LIMIT 1`
	err := cockroach.Select(ctx, sqlQuery, ID).ScanOne(deviceProcessStatusHistoryData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return deviceProcessStatusHistoryData, nil
}

func (i *sDeviceBrokenHistoryRepo) Insert(ctx context.Context, e *model.DeviceBrokenHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (i *sDeviceBrokenHistoryRepo) Update(ctx context.Context, e *model.DeviceBrokenHistory) error {
	return cockroach.Update(ctx, e)
}

// SearchDeviceBrokenHistoryOpts all params is options
type SearchDeviceBrokenHistoryOpts struct {
	IDs           []string
	ProcessStatus []int8
	CreatedFrom   time.Time
	CreatedTo     time.Time
	DeviceID      string
	DeviceIDs     []string
	IsResolved    int16
	ErrorCodes    []string
	//ProductionOrderStageID string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchDeviceBrokenHistoryOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.DeviceBrokenHistoryFieldID)
	}

	if len(s.ProcessStatus) > 0 {
		args = append(args, s.ProcessStatus)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceBrokenHistoryFieldProcessStatus, len(args))
	}

	if s.IsResolved != 0 {
		args = append(args, s.IsResolved)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceBrokenHistoryFieldIsResolved, len(args))
	}

	if len(s.ErrorCodes) > 0 {
		args = append(args, s.ErrorCodes)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceBrokenHistoryFieldErrorCode, len(args))
	}
	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceBrokenHistoryFieldDeviceID, len(args))
	}

	if len(s.DeviceIDs) > 0 {
		args = append(args, s.DeviceIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceBrokenHistoryFieldDeviceID, len(args))
	}

	//if s.ProductionOrderStageID != "" {
	//	args = append(args, s.ProductionOrderStageID)
	//	conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceBrokenHistoryFieldProductionOrderStageDeviceID, len(args))
	//}
	if s.CreatedFrom.IsZero() == false {
		args = append(args, s.CreatedFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.DeviceBrokenHistoryFieldCreatedAt, len(args))
	}
	if s.CreatedTo.IsZero() == false {
		args = append(args, s.CreatedTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.DeviceBrokenHistoryFieldCreatedAt, len(args))
	}

	b := &model.DeviceBrokenHistory{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}

	joins += fmt.Sprintf(" LEFT JOIN users AS u ON u.id = b.updated_by ")
	joins += fmt.Sprintf(" LEFT JOIN users AS u2 ON u2.id = b.created_by ")
	joins += fmt.Sprintf(" LEFT JOIN production_order_stage_devices AS posd ON posd.id = b.production_order_stage_device_id ")
	joins += fmt.Sprintf(" LEFT JOIN production_order_stages AS pos ON pos.id = posd.production_order_stage_id ")
	return fmt.Sprintf("SELECT b.%s,u2.name AS created_user_name,u.name AS updated_user_name, pos.stage_id as stage_id  FROM %s AS b %s WHERE TRUE %s %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type DeviceBrokenHistoryData struct {
	*model.DeviceBrokenHistory
	StageID         sql.NullString `db:"stage_id"`
	UpdatedUserName sql.NullString `db:"updated_user_name"`
	CreatedUserName sql.NullString `db:"created_user_name"`
}

func (i *sDeviceBrokenHistoryRepo) Search(ctx context.Context, s *SearchDeviceBrokenHistoryOpts) ([]*DeviceBrokenHistoryData, error) {
	DeviceBrokenHistory := make([]*DeviceBrokenHistoryData, 0)
	sqlQuery, args := s.buildQuery(false)

	err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&DeviceBrokenHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return DeviceBrokenHistory, nil
}

func (i *sDeviceBrokenHistoryRepo) Count(ctx context.Context, s *SearchDeviceBrokenHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sqlQuery, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sqlQuery, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sDeviceBrokenHistoryRepo.Count: %w", err)
	}

	return countResult, nil
}
