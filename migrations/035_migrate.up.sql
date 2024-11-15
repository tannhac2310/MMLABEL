ALTER TABLE production_order_device_config
    ADD COLUMN device_type varchar(250) NOT NULL DEFAULT 'printer'; -- printer, other.

DROP INDEX idx_custom_fields_entity_id_entity_type_field CASCADE;
