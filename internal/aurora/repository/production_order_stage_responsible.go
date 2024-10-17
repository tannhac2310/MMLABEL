package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductionOrderStageResponsibleRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrderStageResponsible) error
	Update(ctx context.Context, e *model.ProductionOrderStageResponsible) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*ProductionOrderStageResponsibleData, error)
	Search(ctx context.Context, s *SearchProductionOrderStageResponsibleOpts) ([]*ProductionOrderStageResponsibleData, error)
	Count(ctx context.Context, s *SearchProductionOrderStageResponsibleOpts) (*CountResult, error)
}

type sProductionOrderStageResponsibleRepo struct {
}

func NewProductionOrderStageResponsibleRepo() ProductionOrderStageResponsibleRepo {
	return &sProductionOrderStageResponsibleRepo{}
}

func (r *sProductionOrderStageResponsibleRepo) Insert(ctx context.Context, e *model.ProductionOrderStageResponsible) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sProductionOrderStageResponsibleRepo) FindByID(ctx context.Context, id string) (*ProductionOrderStageResponsibleData, error) {
	e := &ProductionOrderStageResponsibleData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sProductionOrderStageResponsibleRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sProductionOrderStageResponsibleRepo) Update(ctx context.Context, e *model.ProductionOrderStageResponsible) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sProductionOrderStageResponsibleRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE production_order_stage_responsible SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("production_order_stage_responsible cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sProductionOrderStageResponsibleRepo not found any records to delete")
	}

	return nil
}

// SearchProductionOrderStageResponsibleOpts all params is options
type SearchProductionOrderStageResponsibleOpts struct {
	IDs                []string
	ProductionOrderIDs []string
	UserIDs            []string
	DeviceIDs          []string
	Limit              int64
	Offset             int64
	Sort               *Sort
}

func (s *SearchProductionOrderStageResponsibleOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " JOIN production_order_stage AS pos ON pos.id = b.production_order_stage_id "
	joins += " JOIN users AS u ON u.id = b.user_id "
	joins += " JOIN production_order AS po ON po.id = pos.production_order_id "
	joins += " JOIN production_order_stage_devices AS posd ON posd.production_order_stage_id = pos.id "

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderStageResponsibleFieldID)
	}
	if len(s.ProductionOrderIDs) > 0 {
		args = append(args, s.ProductionOrderIDs)
		conds += fmt.Sprintf(" AND pos.production_order_id = ANY($%d)", len(args))
	}
	if len(s.UserIDs) > 0 {
		args = append(args, s.UserIDs)
		conds += fmt.Sprintf(" AND b.user_id = ANY($%d)", len(args))
	}
	if len(s.DeviceIDs) > 0 {
		args = append(args, s.DeviceIDs)
		conds += fmt.Sprintf(" AND posd.device_id = ANY($%d)", len(args))
	}

	b := &model.ProductionOrderStageResponsible{}
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

type ProductionOrderStageResponsibleData struct {
	*model.ProductionOrderStageResponsible
}

func (r *sProductionOrderStageResponsibleRepo) Search(ctx context.Context, s *SearchProductionOrderStageResponsibleOpts) ([]*ProductionOrderStageResponsibleData, error) {
	ProductionOrderStageResponsible := make([]*ProductionOrderStageResponsibleData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&ProductionOrderStageResponsible)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return ProductionOrderStageResponsible, nil
}

func (r *sProductionOrderStageResponsibleRepo) Count(ctx context.Context, s *SearchProductionOrderStageResponsibleOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sProductionOrderStageResponsibleRepo.Count: %w", err)
	}

	return countResult, nil
}
