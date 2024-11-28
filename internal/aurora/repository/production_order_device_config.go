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

type ProductionOrderDeviceConfigRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrderDeviceConfig) error
	Update(ctx context.Context, e *model.ProductionOrderDeviceConfig) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) ([]*ProductionOrderDeviceConfigData, error)
	Count(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) (*CountResult, error)
}

type sProductionOrderDeviceConfigRepo struct {
}

func NewProductionOrderDeviceConfigRepo() ProductionOrderDeviceConfigRepo {
	return &sProductionOrderDeviceConfigRepo{}
}

func (r *sProductionOrderDeviceConfigRepo) Insert(ctx context.Context, e *model.ProductionOrderDeviceConfig) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductionOrderDeviceConfigRepo) Update(ctx context.Context, e *model.ProductionOrderDeviceConfig) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sProductionOrderDeviceConfigRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE production_order_device_config SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("production_order_device_config cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sProductionOrderDeviceConfigRepo not found any records to delete")
	}

	return nil
}

// SearchProductionOrderDeviceConfigOpts all params is options
type SearchProductionOrderDeviceConfigOpts struct {
	IDs               []string
	Search            string
	ProductionOrderID string
	ProductionPlanID  string
	DeviceType        string
	MasterDataIDS     []string
	InkIDs            []string
	Limit             int64
	Offset            int64
	Sort              *Sort
}

func (s *SearchProductionOrderDeviceConfigOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " Left JOIN production_orders AS po ON po.id = b.production_order_id"
	// join production_plans table
	joins += " Left JOIN production_plans AS pp ON pp.id = b.production_plan_id"
	// join devices table
	joins += " JOIN devices AS d ON d.id = b.device_id"

	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductQualityFieldProductionOrderID, len(args))
	}
	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderDeviceConfigFieldID)
	}
	if s.ProductionPlanID != "" {
		args = append(args, s.ProductionPlanID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderDeviceConfigFieldProductionPlanID, len(args))
	}
	if s.DeviceType != "" {
		args = append(args, s.DeviceType)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderDeviceConfigFieldDeviceType, len(args))
	}

	if len(s.MasterDataIDS) > 0 {
		args = append(args, s.MasterDataIDS)
		// = array ma_phim or ma_khung or ma_mau_muc
		conds += fmt.Sprintf(" AND (b.ma_phim = ANY($%d) OR b.ma_khung = ANY($%d) OR b.ma_mau_muc = ANY($%d))", len(args), len(args), len(args))
	}

	if len(s.InkIDs) > 0 {
		args = append(args, s.InkIDs)
		conds += fmt.Sprintf(" AND b.ink_id = ANY($%d)", len(args))
	}

	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (b.id like $%[1]d  or b.color ILIKE $%[1]d OR d.name ILIKE $%[1]d OR d.code ILIKE $%[1]d OR po.name ILIKE $%[1]d) or pp.name ILIKE $%[1]d", len(args))
	}

	b := &model.ProductionOrderDeviceConfig{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}

	return fmt.Sprintf(`SELECT b.%s, po.name as production_order_name, 
       d.name as device_name, d.code as device_code, pp.name as production_plan_name
FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductionOrderDeviceConfigData struct {
	*model.ProductionOrderDeviceConfig
	ProductionOrderName sql.NullString `db:"production_order_name"`
	DeviceName          sql.NullString `db:"device_name"`
	DeviceCode          sql.NullString `db:"device_code"`
	ProductionPlanName  sql.NullString `db:"production_plan_name"`
}

func (r *sProductionOrderDeviceConfigRepo) Search(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) ([]*ProductionOrderDeviceConfigData, error) {
	ProductionOrderDeviceConfig := make([]*ProductionOrderDeviceConfigData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&ProductionOrderDeviceConfig)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ProductionOrderDeviceConfig, nil
}

func (r *sProductionOrderDeviceConfigRepo) Count(ctx context.Context, s *SearchProductionOrderDeviceConfigOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sProductionOrderDeviceConfigRepo.Count: %w", err)
	}

	return countResult, nil
}
