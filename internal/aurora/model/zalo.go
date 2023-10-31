package model

import (
	"database/sql"
	"time"
)

const (
	ZaloFieldID             = "id"
	ZaloFieldAppID          = "app_id"
	ZaloFieldSecretKey      = "secret_key"
	ZaloFieldOaID           = "oa_id"
	ZaloFieldOaName         = "oa_name"
	ZaloFieldAccessToken    = "access_token"
	ZaloFieldRefreshToken   = "refresh_token"
	ZaloFieldExpiresIn      = "expires_in"
	ZaloFieldResponsibility = "responsibility"
	ZaloFieldCreatedBy      = "created_by"
	ZaloFieldUpdatedBy      = "updated_by"
	ZaloFieldCreatedAt      = "created_at"
	ZaloFieldUpdatedAt      = "updated_at"
	ZaloFieldDeletedAt      = "deleted_at"
)

type ZaloResponsibility struct {
	UserID   string `json:"userId"`
	Priority int    `json:"priority"`
}
type Zalo struct {
	ID             string                `db:"id"`
	AppID          string                `db:"app_id"`
	SecretKey      string                `db:"secret_key"`
	OaID           string                `db:"oa_id"`
	OaName         string                `db:"oa_name"`
	AccessToken    string                `db:"access_token"`
	RefreshToken   string            `db:"refresh_token"`
	ExpiresIn      int                   `db:"expires_in"`
	Responsibility []*ZaloResponsibility `db:"responsibility"`
	CreatedBy      string                `db:"created_by"`
	UpdatedBy      string            `db:"updated_by"`
	CreatedAt      time.Time         `db:"created_at"`
	UpdatedAt      time.Time         `db:"updated_at"`
	DeletedAt      sql.NullTime      `db:"deleted_at"`
}

func (rcv *Zalo) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ZaloFieldID,
		ZaloFieldAppID,
		ZaloFieldSecretKey,
		ZaloFieldOaID,
		ZaloFieldOaName,
		ZaloFieldAccessToken,
		ZaloFieldRefreshToken,
		ZaloFieldExpiresIn,
		ZaloFieldResponsibility,
		ZaloFieldCreatedBy,
		ZaloFieldUpdatedBy,
		ZaloFieldCreatedAt,
		ZaloFieldUpdatedAt,
		ZaloFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.AppID,
		&rcv.SecretKey,
		&rcv.OaID,
		&rcv.OaName,
		&rcv.AccessToken,
		&rcv.RefreshToken,
		&rcv.ExpiresIn,
		&rcv.Responsibility,
		&rcv.CreatedBy,
		&rcv.UpdatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
	}

	return
}

func (*Zalo) TableName() string {
	return "zalo"
}
