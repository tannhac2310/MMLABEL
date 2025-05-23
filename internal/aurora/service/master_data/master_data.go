package master_data

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Service interface {
	CreateMasterData(ctx context.Context, opt *CreateMasterDataOpts) (string, error)
	UpdateMasterData(ctx context.Context, opt *UpdateMasterDataOpts) error
	DeleteMasterData(ctx context.Context, opt *DeleteMasterDataOpts) error
	FindMasterData(ctx context.Context, opt *FindMasterDataOpts) ([]*MasterData, int64, error)
}

type masterDataService struct {
	masterDataRep                   repository.MasterDataRepo
	masterDataUserField             repository.MasterDataUserFieldRepo
	customFieldRepo                 repository.CustomFieldRepo
	productionOrderDeviceConfigRepo repository.ProductionOrderDeviceConfigRepo
}

func NewService(masterDataRep repository.MasterDataRepo, masterDataUserField repository.MasterDataUserFieldRepo,
	customFieldRepo repository.CustomFieldRepo,
	productionOrderDeviceConfigRepo repository.ProductionOrderDeviceConfigRepo,
) Service {
	return &masterDataService{
		masterDataRep:                   masterDataRep,
		masterDataUserField:             masterDataUserField,
		customFieldRepo:                 customFieldRepo,
		productionOrderDeviceConfigRepo: productionOrderDeviceConfigRepo,
	}
}

type CreateMasterDataUserField struct {
	FieldName  string
	FieldValue string
}

type CreateMasterDataOpts struct {
	Type        enum.MasterDataType
	Name        string
	Code        string
	Description string
	Status      enum.MasterDataStatus
	UserFields  []CreateMasterDataUserField
	CreatedBy   string
}

type UpdateMasterDataOpts struct {
	ID          string
	Name        string
	Code        string
	Description string
	UserFields  []CreateMasterDataUserField
	Status      enum.MasterDataStatus
	UpdateBy    string
}

type DeleteMasterDataOpts struct {
	ID string
}

type FindMasterDataOpts struct {
	IDs    []string
	Type   enum.MasterDataType
	Search string
	Limit  int64
	Offset int64
}

type MasterData struct {
	ID                string
	Type              enum.MasterDataType
	Name              string
	Code              string
	Description       string
	Status            enum.MasterDataStatus
	UserFields        []*MasterDataUserField
	ProductionPlanIDs []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CreatedBy         string
	UpdatedBy         string
}

type MasterDataUserField struct {
	ID           string
	MasterDataID string
	FieldName    string
	FieldValue   string
}
