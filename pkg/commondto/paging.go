package commondto

import (
	"github.com/buger/jsonparser"
)

type Paging struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
type Sort struct {
	Order string `json:"order"`
	By    string `json:"by"`
}
type RespPaging struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
	Total  int64 `json:"total"`
}

func (p *Paging) UnmarshalJSON(data []byte) error {
	if value, err := jsonparser.GetInt(data, "limit"); err == nil {
		p.Limit = value
	}

	if value, err := jsonparser.GetInt(data, "offset"); err == nil {
		p.Offset = value
	}

	if p.Limit == 0 {
		p.Limit = 1
	}

	return nil
}
