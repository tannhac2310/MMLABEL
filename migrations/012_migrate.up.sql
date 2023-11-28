ALTER TABLE production_orders ADD COLUMN estimated_start_at TIMESTAMPTZ;
ALTER TABLE production_orders ADD COLUMN estimated_complete_at TIMESTAMPTZ;
