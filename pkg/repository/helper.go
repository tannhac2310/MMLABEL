package repository

import (
	"errors"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func permissionCondition(entity enum.PermissionEntity, userID, targetAlias string) string {
	p := &model.Permission{}
	return fmt.Sprintf(" JOIN %[1]s  ON %[1]s.entity = %[2]d and %[1]s.user_id = '%[3]s' and %[1]s.element_id = %[4]s.id ", p.TableName(), entity, userID, targetAlias)
}

var (
	ErrNotFound = errors.New("not found")
)
