package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InspectionFormRepo interface {
	Insert(ctx context.Context, e *model.InspectionForm) error
	Update(ctx context.Context, e *model.InspectionForm) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*InspectionFormData, error)
	Search(ctx context.Context, s *SearchInspectionFormOpts) ([]*InspectionFormData, error)
	Count(ctx context.Context, s *SearchInspectionFormOpts) (*CountResult, error)
	CountAll(ctx context.Context) (*int64, error)
}

type sInspectionFormRepo struct {
}

func NewInspectionFormRepo() InspectionFormRepo {
	return &sInspectionFormRepo{}
}

func (r *sInspectionFormRepo) CountAll(ctx context.Context) (*int64, error) {
	var count int64
	err := cockroach.Select(ctx, "SELECT count(*) as cnt FROM inspection_forms").ScanOne(&count)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return &count, nil

}

func (r *sInspectionFormRepo) Insert(ctx context.Context, e *model.InspectionForm) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sInspectionFormRepo) FindByID(ctx context.Context, id string) (*InspectionFormData, error) {
	e := &InspectionFormData{
		InspectionForm: &model.InspectionForm{},
	}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sInspectionFormRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sInspectionFormRepo) Update(ctx context.Context, e *model.InspectionForm) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sInspectionFormRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE inspection_forms SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("inspection_forms cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sInspectionFormRepo not found any records to delete")
	}

	return nil
}

// SearchInspectionFormOpts all params is options
type SearchInspectionFormOpts struct {
	IDs               []string
	ProductionOrderID string
	DefectType        []string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	ProductSearch     string
	CustomerSearch    string
	Limit             int64
	Offset            int64
	Sort              *Sort
}

func (s *SearchInspectionFormOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""
	joins += " JOIN production_orders AS po ON b.production_order_id = po.id"
	joins += " JOIN products AS p ON p.id = b.product_id "
	joins += " JOIN customers AS c ON c.id = p.customer_id "
	joins += " JOIN orders AS o ON o.id = po.order_id "
	joins += " LEFT JOIN users AS u ON u.id = b.created_by "
	if s.ProductSearch != "" {
		args = append(args, "%"+s.ProductSearch+"%")
		conds += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.code ILIKE $%d)", len(args), len(args))
	}
	if s.CustomerSearch != "" {
		// search like name and code
		args = append(args, "%"+s.CustomerSearch+"%")
		conds += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.code ILIKE $%d)", len(args), len(args))
	}

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.InspectionFormFieldID)
	}
	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InspectionFormFieldProductionOrderID, len(args))
	}
	if len(s.DefectType) > 0 {
		args = append(args, s.DefectType)
		joins += " JOIN inspection_errors AS ie ON b.id = ie.inspection_form_id "
		conds += fmt.Sprintf(" AND ie.%s = ANY($%d)", model.InspectionErrorFieldErrorType, len(args))
	}

	if !s.CreatedAtFrom.IsZero() {
		args = append(args, s.CreatedAtFrom)
		conds += fmt.Sprintf(" AND b.inspection_date >= $%d", len(args))
	}
	if !s.CreatedAtTo.IsZero() {
		args = append(args, s.CreatedAtTo)
		conds += fmt.Sprintf(" AND b.inspection_date <= $%d", len(args))
	}

	b := &model.InspectionForm{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.created_at DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`
SELECT b.%s, po.product_name as production_order_name,  po.product_code as production_order_code, po.order_id as order_id,
p.code as product_code, p.name as product_name,
    c.id as customer_id,
c.name as customer_name, c.code as customer_code,
-- po.order_id as ma_don_dat_hang,
    o.ma_dat_hang_mm as ma_don_dat_hang,
    o.status as trang_thai_don_hang,
    u.name as created_by_name
FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d
`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type InspectionFormData struct {
	*model.InspectionForm
	OrderID             string `db:"order_id"`
	ProductionOrderName string `db:"production_order_name"`
	ProductionOrderCode string `db:"production_order_code"`
	ProductCode         string `db:"product_code"`
	ProductID           string `db:"product_id"`
	ProductName         string `db:"product_name"`
	CustomerID          string `db:"customer_id"`
	CustomerName        string `db:"customer_name"`
	CustomerCode        string `db:"customer_code"`
	CreatedByName       string `db:"created_by_name"`
	MaDonDatHang        string `db:"ma_don_dat_hang"`
	TrangThaiDonHang    string `db:"trang_thai_don_hang"`
}

func (r *sInspectionFormRepo) Search(ctx context.Context, s *SearchInspectionFormOpts) ([]*InspectionFormData, error) {
	v := make([]*InspectionFormData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&v)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return v, nil
}

func (r *sInspectionFormRepo) Count(ctx context.Context, s *SearchInspectionFormOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sInspectionFormRepo.Count: %w", err)
	}

	return countResult, nil
}
