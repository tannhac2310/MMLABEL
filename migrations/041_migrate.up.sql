alter table ink
    add column kho varchar(255) default '',
    add column loai_muc varchar(255) default '',
    add column nha_cung_cap varchar(255) default '',
    add column tinh_trang varchar(255) default '';

ALTER TABLE orders
    ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE order_items
    ADD COLUMN deleted_at TIMESTAMPTZ;