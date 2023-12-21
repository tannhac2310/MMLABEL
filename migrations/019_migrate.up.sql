-- add column device_ids into product_quantity table
ALTER TABLE product_quality ADD device_ids STRING[]  NOT NULL DEFAULT ARRAY[];