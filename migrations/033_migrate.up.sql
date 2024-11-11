ALTER TABLE customers
    ADD COLUMN data jsonb;
ALTER TABLE customers
    ADD COLUMN search_content text;

ALTER TABLE master_data
    ADD COLUMN code varchar(255) NOT NULL default '';