-- create table production_order_stage_responsible

CREATE TABLE production_order_stage_responsible
(
    id                               VARCHAR(50) NOT NULL,
    po_stage_device_id varchar(50) NOT NULL,
    user_id                          varchar(50) NOT NULL,
    deleted_at                       TIMESTAMPTZ,
    CONSTRAINT pk_production_order_stage_responsible PRIMARY KEY (id ASC)
);
create index idx_posp_pos_id on production_order_stage_responsible(po_stage_device_id);
create index idx_posp_user_id on production_order_stage_responsible(user_id);


ALTER TABLE production_order_stage_devices
    ADD COLUMN estimated_start_at TIMESTAMP;

ALTER TABLE production_order_stage_devices
    ADD COLUMN color VARCHAR(250);

ALTER TABLE production_order_stage_devices
    ADD COLUMN po_device_config_id VARCHAR(250);

ALTER TABLE production_order_stage_devices
    ADD COLUMN start_at TIMESTAMP;

ALTER TABLE production_order_stage_devices
    ADD COLUMN complete_at TIMESTAMP;

-- write sql migrate production_order_stage_devices.responsible to production_order_stage_responsible
WITH production_order_stage_devices AS (
    SELECT
        id,
        unnest(responsible) AS user_id
    FROM production_order_stage_devices
    WHERE responsible IS NOT NULL
)
INSERT INTO production_order_stage_responsible (id, po_stage_device_id, user_id)
SELECT
    uuid_generate_v4(),
    id,
    user_id
FROM production_order_stage_devices;