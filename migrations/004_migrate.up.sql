ALTER TABLE production_order_stages
    ADD estimated_start_at TIMESTAMPTZ;
ALTER TABLE production_order_stages
    ADD estimated_complete_at TIMESTAMPTZ;