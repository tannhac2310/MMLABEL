-- add column time_spend into event_logs table
ALTER TABLE event_logs ADD time_spend INT DEFAULT 0;

ALTER TABLE device_working_history ALTER COLUMN production_order_stage_device_id DROP NOT NULL;

ALTER TABLE device_working_history ADD po_quantity INT DEFAULT 0;
ALTER TABLE device_working_history ADD po_working_time INT DEFAULT 0;

