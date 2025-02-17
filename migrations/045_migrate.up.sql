
ALTER TABLE production_orders
    ADD COLUMN order_id varchar(255) default NULL;
ALTER TABLE production_order_stages
    ADD COLUMN ghi_chu_ban_in_nguon TEXT default NULL;