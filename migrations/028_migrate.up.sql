alter table production_order_device_config add column production_plan_id varchar(255) null;
    ALTER TABLE production_order_device_config
        ALTER COLUMN production_order_id DROP NOT NULL;