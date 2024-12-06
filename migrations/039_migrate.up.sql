drop table if exists orders;

CREATE TABLE orders
(
    id                     VARCHAR(50) PRIMARY KEY,
    title                  VARCHAR(255) NOT NULL,
    ma_dat_hang_mm         VARCHAR(255) NOT NULL,
    ma_hop_dong_khach_hang VARCHAR(255) NOT NULL,
    ma_hop_dong            VARCHAR(255) NOT NULL,
    sale_name              VARCHAR(255),
    sale_admin_name        VARCHAR(255),
    status                 VARCHAR(50)  NOT NULL,
    created_by             VARCHAR(50)  NOT NULL,
    updated_by             VARCHAR(50)  NOT NULL,
    created_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- add column to ink
ALTER TABLE ink
    ADD COLUMN kho VARCHAR(255),
    ADD COLUMN loai_muc VARCHAR(255),
    ADD COLUMN nha_cung_cap  VARCHAR(255),
    ADD COLUMN tinh_trang VARCHAR(255);

