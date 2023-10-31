package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

type CasbinRuleRepo interface {
	Insert(ctx context.Context, e *model.CasbinRule) error
	Delete(ctx context.Context, v0, v1 string) error
	Search(ctx context.Context, s *SearchCasbinRulesOpts) ([]*CasbinRuleData, error)
	Count(ctx context.Context, s *SearchCasbinRulesOpts) (*CountResult, error)
}

type casbinRuleRepo struct {
}
type CasbinRuleData struct {
	*model.CasbinRule
}

func (r *casbinRuleRepo) Count(ctx context.Context, s *SearchCasbinRulesOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("stage.Count: %w", err)
	}

	return countResult, nil
}

func NewCasbinRuleRepo() CasbinRuleRepo {
	return &casbinRuleRepo{}
}

func (r *casbinRuleRepo) Insert(ctx context.Context, e *model.CasbinRule) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *casbinRuleRepo) Delete(ctx context.Context, v0, v1 string) error {
	sql := `DELETE FROM casbin_rules	WHERE v0 = $1 and v1 = $2`

	cmd, err := cockroach.Exec(ctx, sql, v0, v1)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

// SearchCasbinRulesOpts all params is options
type SearchCasbinRulesOpts struct {
	IDs    []string
	Name   string
	Limit  int64
	Offset int64
}

func (s *SearchCasbinRulesOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""

	b := &model.CasbinRule{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(1) as cnt
		FROM %s AS b
		WHERE TRUE %s`, b.TableName(), conds), args
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b 
		WHERE TRUE %s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), conds, s.Limit, s.Offset), args
}

func (r *casbinRuleRepo) Search(ctx context.Context, s *SearchCasbinRulesOpts) ([]*CasbinRuleData, error) {
	casbinRules := make([]*CasbinRuleData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&casbinRules)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return casbinRules, nil
}
