package model

const (
	CasbinRuleFieldPType = "p_type"
	CasbinRuleFieldV0    = "v0"
	CasbinRuleFieldV1    = "v1"
	CasbinRuleFieldV2    = "v2"
	CasbinRuleFieldV3    = "v3"
	CasbinRuleFieldV4    = "v4"
	CasbinRuleFieldV5    = "v5"
	CasbinRuleFieldRowid = "rowid"
)

type CasbinRule struct {
	PType string `db:"p_type"`
	V0    string `db:"v0"`
	V1    string `db:"v1"`
	V2    string `db:"v2"`
	V3    string `db:"v3"`
	V4    string `db:"v4"`
	V5    string `db:"v5"`
	Rowid int64  `db:"rowid"`
}

func (rcv *CasbinRule) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		CasbinRuleFieldPType,
		CasbinRuleFieldV0,
		CasbinRuleFieldV1,
		CasbinRuleFieldV2,
		CasbinRuleFieldV3,
		CasbinRuleFieldV4,
		CasbinRuleFieldV5,
		CasbinRuleFieldRowid,
	}

	values = []interface{}{
		&rcv.PType,
		&rcv.V0,
		&rcv.V1,
		&rcv.V2,
		&rcv.V3,
		&rcv.V4,
		&rcv.V5,
		&rcv.Rowid,
	}

	return
}

func (*CasbinRule) TableName() string {
	return "casbin_rules"
}
