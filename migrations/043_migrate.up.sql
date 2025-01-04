ALTER TABLE custom_fields
    ADD COLUMN data jsonb default '{}';
ALTER TABLE production_orders
    ADD COLUMN data jsonb default '{}';
ALTER TABLE production_orders
    ADD COLUMN workflows jsonb default '[]';