package repository

import (
	"context"
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
	IDs    []string
	Search string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchProductionOrderDeviceConfigOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " JOIN production_orders AS po ON po.id = b.production_order_id"

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderDeviceConfigFieldID)
	}
	// todo add more search options example:
	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d OR po.name ILIKE $%[1]d)",
			len(args), model.ProductionOrderDeviceConfigFieldSearch, model.ProductionOrderDeviceConfigFieldColor)
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

	return fmt.Sprintf(`SELECT b.%s, po.name as production_order_name 
FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductionOrderDeviceConfigData struct {
	*model.ProductionOrderDeviceConfig
	ProductionOrderName string `db:"production_order_name"`
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
