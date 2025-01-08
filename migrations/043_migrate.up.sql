ALTER TABLE custom_fields
    ADD COLUMN data jsonb default '{}';
ALTER TABLE production_orders
    ADD COLUMN data jsonb default '{}';
ALTER TABLE production_orders
    ADD COLUMN workflows jsonb default '[]';
ALTER TABLE production_order_stages
    ADD COLUMN so_luong integer default 0;