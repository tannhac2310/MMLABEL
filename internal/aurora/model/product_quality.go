package model

// TODO remove this file
//import (
//	"database/sql"
//	"time"
//)
//
//const (
//	ProductQualityFieldID                = "id"
//	ProductQualityFieldProductionOrderID = "production_order_id"
//	ProductQualityFieldProductID         = "product_id"
//	ProductQualityFieldDefectType        = "defect_type"
//	ProductQualityFieldDefectCode        = "defect_code"
//	ProductQualityFieldDefectLevel       = "defect_level"
//	ProductQualityFieldProductionStageID = "production_stage_id"
//	ProductQualityFieldDefectiveQuantity = "defective_quantity"
//	ProductQualityFieldGoodQuantity      = "good_quantity"
//	ProductQualityFieldDescription       = "description"
//	ProductQualityFieldCreatedBy         = "created_by"
//	ProductQualityFieldCreatedAt         = "created_at"
//	ProductQualityFieldUpdatedAt         = "updated_at"
//	ProductQualityFieldDeletedAt         = "deleted_at"
//	ProductQualityFieldDeviceIDs         = "device_ids"
//)
//
//type ProductQuality struct {
//	ID                string         `db:"id"`
//	ProductionOrderID sql.NullString `db:"production_order_id"`
//	ProductID         sql.NullString `db:"product_id"`
//	DefectTypes        sql.NullString `db:"defect_type"`
//	DefectCode        sql.NullString `db:"defect_code"`
//	DefectLevel       int16          `db:"defect_level"`
//	ProductionStageID sql.NullString `db:"production_stage_id"`
//	DefectiveQuantity int64          `db:"defective_quantity"`
//	GoodQuantity      int64          `db:"good_quantity"`
//	Description       sql.NullString `db:"description"`
//	CreatedBy         string         `db:"created_by"`
//	CreatedAt         time.Time      `db:"created_at"`
//	UpdatedAt         time.Time      `db:"updated_at"`
//	DeletedAt         sql.NullTime   `db:"deleted_at"`
//	DeviceIDs         []string       `db:"device_ids"`
//}
//
//func (rcv *ProductQuality) FieldMap() (fields []string, values []interface{}) {
//	fields = []string{
//		ProductQualityFieldID,
//		ProductQualityFieldProductionOrderID,
//		ProductQualityFieldProductID,
//		ProductQualityFieldDefectType,
//		ProductQualityFieldDefectCode,
//		ProductQualityFieldDefectLevel,
//		ProductQualityFieldProductionStageID,
//		ProductQualityFieldDefectiveQuantity,
//		ProductQualityFieldGoodQuantity,
//		ProductQualityFieldDescription,
//		ProductQualityFieldCreatedBy,
//		ProductQualityFieldCreatedAt,
//		ProductQualityFieldUpdatedAt,
//		ProductQualityFieldDeletedAt,
//		ProductQualityFieldDeviceIDs,
//	}
//
//	values = []interface{}{
//		&rcv.ID,
//		&rcv.ProductionOrderID,
//		&rcv.ProductID,
//		&rcv.DefectTypes,
//		&rcv.DefectCode,
//		&rcv.DefectLevel,
//		&rcv.ProductionStageID,
//		&rcv.DefectiveQuantity,
//		&rcv.GoodQuantity,
//		&rcv.Description,
//		&rcv.CreatedBy,
//		&rcv.CreatedAt,
//		&rcv.UpdatedAt,
//		&rcv.DeletedAt,
//		&rcv.DeviceIDs,
//	}
//
//	return
//}
//
//func (*ProductQuality) TableName() string {
//	return "product_quality"
//}
