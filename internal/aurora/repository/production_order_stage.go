package repository

import (
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductionOrderStageRepo interface {
	Insert(ctx context.Context, e *model.ProductionOrderStage) error
	Update(ctx context.Context, e *model.ProductionOrderStage) error
	FindByID(ctx context.Context, id string) (*model.ProductionOrderStage, error)
	SoftDelete(ctx context.Context, id string) error
	SoftDeletes(ctx context.Context, ids []string) error
	Search(ctx context.Context, s *SearchProductionOrderStagesOpts) ([]*model.ProductionOrderStage, error)
	Count(ctx context.Context, s *SearchProductionOrderStagesOpts) (*CountResult, error)
	DeleteByProductionOrderID(ctx context.Context, id string) error
}

type productionOrderStagesRepo struct {
}

func NewProductionOrderStageRepo() ProductionOrderStageRepo {
	return &productionOrderStagesRepo{}
}
func (p *productionOrderStagesRepo) FindByID(ctx context.Context, id string) (*model.ProductionOrderStage, error) {
	e := &model.ProductionOrderStage{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("cockroach.FindOne: %w", err)
	}

	return e, nil
}

func (p *productionOrderStagesRepo) DeleteByProductionOrderID(ctx context.Context, id string) error {
	sql := `UPDATE production_order_stages
		SET deleted_at = NOW()
		WHERE production_order_id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *productionOrderStagesRepo) Insert(ctx context.Context, e *model.ProductionOrderStage) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *productionOrderStagesRepo) Update(ctx context.Context, e *model.ProductionOrderStage) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *productionOrderStagesRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE production_order_stages
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
func (r *productionOrderStagesRepo) SoftDeletes(ctx context.Context, ids []string) error {
	sql := `UPDATE production_order_stages
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

// SearchProductionOrderStagesOpts all params is options
type SearchProductionOrderStagesOpts struct {
	IDs                        []string
	ProductionOrderID          string
	ProductionOrderStageStatus enum.ProductionOrderStageStatus
	Limit                      int64
	Offset                     int64
	Sort                       *Sort
}

func (s *SearchProductionOrderStagesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.ProductionOrderStageFieldID)
	}
	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageFieldProductionOrderID, len(args))
	}
	if s.ProductionOrderStageStatus > 0 {
		args = append(args, s.ProductionOrderStageStatus)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductionOrderStageFieldStatus, len(args))
	}

	b := &model.ProductionOrderStage{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.sorting DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

func (r *productionOrderStagesRepo) Search(ctx context.Context, s *SearchProductionOrderStagesOpts) ([]*model.ProductionOrderStage, error) {
	message := make([]*model.ProductionOrderStage, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *productionOrderStagesRepo) Count(ctx context.Context, s *SearchProductionOrderStagesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
