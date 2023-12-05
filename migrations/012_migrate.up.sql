ALTER TABLE production_orders ADD COLUMN estimated_start_at TIMESTAMPTZ;
ALTER TABLE production_orders ADD COLUMN estimated_complete_at TIMESTAMPTZ;


-- DROP INDEX  idx_production_order_name;
CREATE UNIQUE INDEX idx_production_order_name2 ON production_orders (name,product_code);
CREATE INDEX idx_production_order_status_2 ON production_orders (status) where deleted_at is null;
CREATE INDEX idx_production_order_stages_production_order_id_2 ON production_order_stages (production_order_id) where deleted_at is null;
CREATE INDEX idx_production_order_stage_devices_ref_id_2 ON production_order_stage_devices (production_order_stage_id) where deleted_at is null;
CREATE INDEX idx_custom_fields_entity ON custom_fields (entity_id,entity_type);
CREATE INDEX idx_production_order_stage_devices_device_id ON production_order_stage_devices (device_id);
CREATE INDEX idx_production_order_stage_devices_2 ON production_order_stage_devices (device_id,production_order_stage_id) where deleted_at is null;
