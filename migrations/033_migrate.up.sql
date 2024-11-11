ALTER TABLE customers
    ADD COLUMN data jsonb;
ALTER TABLE customers
    ADD COLUMN search_content text;