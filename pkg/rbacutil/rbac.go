package rbacutil

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
)

const authzModel = `[request_definition]
r = sub, obj

[policy_definition]
p = sub, obj

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && g2(r.obj, p.obj) && keyMatch(r.obj, p.obj) || r.sub == "root"
`

func New(uri string) (casbin.IEnforcer, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	m, err := model.NewModelFromString(authzModel)
	if err != nil {
		return nil, err
	}

	tableName := "casbin_rules"
	a, err := casbinpgadapter.NewAdapter(db, tableName)
	if err != nil {
		return nil, err
	}

	e, _ := casbin.NewSyncedEnforcer(m, a)
	e.StartAutoLoadPolicy(10 * time.Second)

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return e, nil
}
