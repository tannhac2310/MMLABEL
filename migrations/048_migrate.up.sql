ALTER TABLE orders
    ADD COLUMN payment_method varchar(255) default '',
    ADD COLUMN payment_method_other varchar(255) default '',
    ADD COLUMN customer_id varchar(50) default '',
    ADD COLUMN customer_address_options text default '',
    ADD COLUMN delivery_address text default '';