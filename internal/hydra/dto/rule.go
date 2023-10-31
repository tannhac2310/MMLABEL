package dto

import "mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"

type Rule struct {
	Role string `json:"role"`
	Rule string `json:"rule"`
}

type CreateRuleRequest struct {
	Role string `json:"role" binding:"required"`
	Rule string `json:"rule" binding:"required"`
}
type CreateRuleResponse struct {
}

type DeleteRuleRequest struct {
	Role string `json:"role" binding:"required"`
	Rule string `json:"rule" binding:"required"`
}
type DeleteRuleResponse struct {
}

type FindRulesRequest struct {
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindRulesResponse struct {
	Rules    []*Rule           `json:"rules"`
	NextPage *commondto.Paging `json:"nextPage"`
	Total    int64             `json:"total"`
}
