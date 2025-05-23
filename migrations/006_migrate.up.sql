CREATE TABLE custom_fields
(
    id          VARCHAR(50) NOT NULL,
    entity_id   VARCHAR(50) NOT NULL,
    entity_type SMALLINT    NOT NULL,
    field       VARCHAR(50) NOT NULL,
    value       TEXT        not null default '',
    description TEXT,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT pk_custom_fields PRIMARY KEY (id ASC)
);
-- Add a unique index to prevent duplicate custom fields
CREATE UNIQUE INDEX idx_custom_fields_entity_id_entity_type_field ON custom_fields (entity_id, entity_type, field);
-- Add a unique index to prevent duplicate name production order
CREATE UNIQUE INDEX idx_production_order_name ON production_orders (name);
