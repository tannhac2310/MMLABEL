package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type MasterDataRepo interface {
	Insert(ctx context.Context, e *model.MasterData) error
	Update(ctx context.Context, e *model.MasterData) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MasterDataData, error)
	Search(ctx context.Context, s *SearchMasterDataOpts) ([]*MasterDataData, error)
	Count(ctx context.Context, s *SearchMasterDataOpts) (*CountResult, error)
}

type sMasterDataRepo struct {
}

func NewMasterDataRepo() MasterDataRepo {
	return &sMasterDataRepo{}
}

func (r *sMasterDataRepo) Insert(ctx context.Context, e *model.MasterData) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMasterDataRepo) FindByID(ctx context.Context, id string) (*MasterDataData, error) {
	e := &MasterDataData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMasterDataRepo) Update(ctx context.Context, e *model.MasterData) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMasterDataRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE master_data SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("master_data cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMasterDataRepo not found any records to delete")
	}

	return nil
}

// SearchMasterDataOpts all params is options
type SearchMasterDataOpts struct {
	IDs          []string
	Type         enum.MasterDataType
	Code         string
	IsIncludeDel bool
	Search       string
	Limit        int64
	Offset       int64
	Sort         *Sort
}

func (s *SearchMasterDataOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MasterDataFieldID)
	}

	if s.Type != "" {
		args = append(args, s.Type)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.MasterDataFieldType, len(args))
	}

	if s.Search != "" {
		args = append(args, "%"+s.Search+"%")
		args = append(args, s.Search)
		conds += fmt.Sprintf(" AND (b.name ILIKE $%d OR b.description ILIKE $%d OR b.code ILIKE $%d OR b.id = $%d)", len(args)-1, len(args)-1, len(args)-1, len(args))
	}

	if s.Code != "" {
		args = append(args, s.Code)
		conds += fmt.Sprintf(" AND b.code = $%d", len(args))
	}

	if !s.IsIncludeDel {
		conds += " AND b.deleted_at IS NULL "
	}

	b := &model.MasterData{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.created_at DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type MasterDataData struct {
	*model.MasterData
}

func (r *sMasterDataRepo) Search(ctx context.Context, s *SearchMasterDataOpts) ([]*MasterDataData, error) {
	MasterData := make([]*MasterDataData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MasterData)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MasterData, nil
}

func (r *sMasterDataRepo) Count(ctx context.Context, s *SearchMasterDataOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataRepo.Count: %w", err)
	}

	return countResult, nil
}
