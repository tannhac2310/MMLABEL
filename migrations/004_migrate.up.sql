ALTER TABLE production_order_stage
    ADD estimated_start_at TIMESTAMPTZ;
ALTER TABLE production_order_stage
    ADD estimated_complete_at TIMESTAMPTZ;