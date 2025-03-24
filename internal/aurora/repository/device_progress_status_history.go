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

type DeviceProgressStatusHistoryRepo interface {
	Insert(ctx context.Context, e *model.DeviceProgressStatusHistory) error
	Update(ctx context.Context, e *model.DeviceProgressStatusHistory) error
	Search(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) ([]*DeviceProgressStatusHistoryData, error)
	Count(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) (*CountResult, error)
	FindProductionOrderStageDeviceID(ctx context.Context, ProductionOrderStageID string, deviceID string) (*DeviceProgressStatusHistoryData, error)
	FindByID(ctx context.Context, ID string) (*DeviceProgressStatusHistoryData, error)
}

type sDeviceProgressStatusHistoryRepo struct {
}

func NewDeviceProgressStatusHistoryRepo() DeviceProgressStatusHistoryRepo {
	return &sDeviceProgressStatusHistoryRepo{}
}
func (i *sDeviceProgressStatusHistoryRepo) FindProductionOrderStageDeviceID(ctx context.Context, ProductionOrderStageDeviceID string, deviceID string) (*DeviceProgressStatusHistoryData, error) {
	deviceProcessStatusHistoryData := &DeviceProgressStatusHistoryData{}
	sqlQuery := `SELECT * FROM device_progress_status_history WHERE production_order_stage_device_id = $1 AND device_id = $2 AND is_resolved = 0 ORDER BY ID DESC LIMIT 1`
	err := cockroach.Select(ctx, sqlQuery, ProductionOrderStageDeviceID, deviceID).ScanOne(deviceProcessStatusHistoryData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return deviceProcessStatusHistoryData, nil
}
func (i *sDeviceProgressStatusHistoryRepo) FindByID(ctx context.Context, ID string) (*DeviceProgressStatusHistoryData, error) {
	deviceProcessStatusHistoryData := &DeviceProgressStatusHistoryData{}
	sqlQuery := `SELECT * FROM device_progress_status_history WHERE ID = $1 ORDER BY ID DESC LIMIT 1`
	err := cockroach.Select(ctx, sqlQuery, ID).ScanOne(deviceProcessStatusHistoryData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	return deviceProcessStatusHistoryData, nil
}
func (i *sDeviceProgressStatusHistoryRepo) Insert(ctx context.Context, e *model.DeviceProgressStatusHistory) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (i *sDeviceProgressStatusHistoryRepo) Update(ctx context.Context, e *model.DeviceProgressStatusHistory) error {
	return cockroach.Update(ctx, e)
}

// SearchDeviceProgressStatusHistoryOpts all params is options
type SearchDeviceProgressStatusHistoryOpts struct {
	IDs           []string
	ProcessStatus []int8
	CreatedFrom   time.Time
	CreatedTo     time.Time
	DeviceID      string
	DeviceIDs     []string
	IsResolved    int16
	ErrorCodes    []string
	//ProductionOrderStageID string
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

	if len(s.ProcessStatus) > 0 {
		args = append(args, s.ProcessStatus)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceProgressStatusHistoryFieldProcessStatus, len(args))
	}

	if s.IsResolved != 0 {
		args = append(args, s.IsResolved)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceProgressStatusHistoryFieldIsResolved, len(args))
	}

	if len(s.ErrorCodes) > 0 {
		args = append(args, s.ErrorCodes)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceProgressStatusHistoryFieldErrorCode, len(args))
	}
	if s.DeviceID != "" {
		args = append(args, s.DeviceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceProgressStatusHistoryFieldDeviceID, len(args))
	}

	if len(s.DeviceIDs) > 0 {
		args = append(args, s.DeviceIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.DeviceProgressStatusHistoryFieldDeviceID, len(args))
	}

	//if s.ProductionOrderStageID != "" {
	//	args = append(args, s.ProductionOrderStageID)
	//	conds += fmt.Sprintf(" AND b.%s = $%d", model.DeviceProgressStatusHistoryFieldProductionOrderStageDeviceID, len(args))
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

type DeviceProgressStatusHistoryData struct {
	*model.DeviceProgressStatusHistory
	StageID         sql.NullString `db:"stage_id"`
	UpdatedUserName sql.NullString `db:"updated_user_name"`
	CreatedUserName sql.NullString `db:"created_user_name"`
}

func (i *sDeviceProgressStatusHistoryRepo) Search(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) ([]*DeviceProgressStatusHistoryData, error) {
	DeviceProgressStatusHistory := make([]*DeviceProgressStatusHistoryData, 0)
	sqlQuery, args := s.buildQuery(false)

	err := cockroach.Select(ctx, sqlQuery, args...).ScanAll(&DeviceProgressStatusHistory)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return DeviceProgressStatusHistory, nil
}

func (i *sDeviceProgressStatusHistoryRepo) Count(ctx context.Context, s *SearchDeviceProgressStatusHistoryOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sqlQuery, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sqlQuery, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sDeviceProgressStatusHistoryRepo.Count: %w", err)
	}

	return countResult, nil
}
