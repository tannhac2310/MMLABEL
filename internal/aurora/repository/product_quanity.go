package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type ProductQualityRepo interface {
	Insert(ctx context.Context, e *model.ProductQuality) error
	Update(ctx context.Context, e *model.ProductQuality) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchProductQualitysOpts) ([]*ProductQualityData, error)
	Analysis(ctx context.Context, s *SearchProductQualitysOpts) ([]*ProductQualityAnalysis, error)
	Count(ctx context.Context, s *SearchProductQualitysOpts) (*CountResult, error)
}

type productQualitysRepo struct {
}

func NewProductQualityRepo() ProductQualityRepo {
	return &productQualitysRepo{}
}

func (r *productQualitysRepo) Insert(ctx context.Context, e *model.ProductQuality) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *productQualitysRepo) Update(ctx context.Context, e *model.ProductQuality) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *productQualitysRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE product_quality
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

// SearchProductQualitysOpts all params is options
type SearchProductQualitysOpts struct {
	ProductionOrderID string
	DefectType        string
	DefectCode        string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	Limit             int64
	Offset            int64
	Sort              *Sort
}

func (s *SearchProductQualitysOpts) buildQuery(isCount bool, isAnalysis bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductQualityFieldProductionOrderID, len(args))
	}
	if s.DefectType != "" {
		args = append(args, s.DefectType)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductQualityFieldDefectType, len(args))
	}
	if s.DefectCode != "" {
		args = append(args, s.DefectCode)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.ProductQualityFieldDefectCode, len(args))
	}

	if s.CreatedAtFrom.IsZero() == false {
		args = append(args, s.CreatedAtFrom)
		conds += fmt.Sprintf(" AND b.%s >= $%d", model.ProductQualityFieldCreatedAt, len(args))
	}
	if s.CreatedAtTo.IsZero() == false {
		args = append(args, s.CreatedAtTo)
		conds += fmt.Sprintf(" AND b.%s <= $%d", model.ProductQualityFieldCreatedAt, len(args))
	}
	b := &model.ProductQuality{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
	}

	if isAnalysis {
		return fmt.Sprintf(`SELECT b.%s, sum(b.defective_quantity) as count
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		GROUP BY b.%s`, model.ProductQualityFieldDefectType, b.TableName(), joins, conds, model.ProductQualityFieldDefectType), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s, po.name as production_order_name
		FROM %s AS b %s
 		JOIN production_orders AS po ON po.id = b.production_order_id
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type ProductQualityData struct {
	*model.ProductQuality
	ProductionOrderName string `db:"production_order_name"`
}

func (r *productQualitysRepo) Search(ctx context.Context, s *SearchProductQualitysOpts) ([]*ProductQualityData, error) {
	productQualitys := make([]*ProductQualityData, 0)
	sql, args := s.buildQuery(false, false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&productQualitys)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return productQualitys, nil
}

type ProductQualityAnalysis struct {
	DefectType string `json:"defectType"`
	Count      int64  `json:"count"`
}

func (r *productQualitysRepo) Analysis(ctx context.Context, s *SearchProductQualitysOpts) ([]*ProductQualityAnalysis, error) {
	productQualitys := make([]*ProductQualityAnalysis, 0)
	sql, args := s.buildQuery(false, true)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&productQualitys)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return productQualitys, nil
}

func (r *productQualitysRepo) Count(ctx context.Context, s *SearchProductQualitysOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true, false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
