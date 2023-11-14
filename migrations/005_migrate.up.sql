ALTER TABLE production_orders
    ADD name VARCHAR(255) NOT NULL;

ALTER TABLE production_order_stages ADD sorting INT NOT NULL DEFAULT 0;