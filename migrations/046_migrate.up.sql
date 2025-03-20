-- Bảng lưu thông tin phiếu kiểm tra sản phẩm
CREATE TABLE inspection_forms
(
    id                    VARCHAR(50) PRIMARY KEY,
    production_order_id     VARCHAR(255)      NOT NULL,
    inspection_date         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    inspector_name          VARCHAR(255) NOT NULL,
    quantity                INTEGER      NOT NULL,
    ma_san_pham             VARCHAR(255) NOT NULL,
    ten_san_pham            VARCHAR(255) NOT NULL,
    so_luong_hop_dong       INTEGER      NOT NULL,
    so_luong_in             INTEGER      NOT NULL,
    ma_don_dat_hang         VARCHAR(255) NOT NULL,
    nguoi_kiem_tra          VARCHAR(255) NOT NULL,
    nguoi_phe_duyet         VARCHAR(255) NOT NULL,
    so_luong_thanh_pham_dat INTEGER      NOT NULL,
    note                    TEXT,
    created_by              VARCHAR(50)  NOT NULL,
    updated_by              VARCHAR(50)  NOT NULL,
    created_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
);

-- Bảng lưu thông tin các lỗi liên quan đến từng phiếu kiểm tra
CREATE TABLE inspection_errors
(
    id                    VARCHAR(50) PRIMARY KEY,
    device_id           VARCHAR(50)  NOT NULL,
    device_name         VARCHAR(255) NOT NULL,
    inspection_form_id VARCHAR(50)  NOT NULL,
    error_type         VARCHAR(255) NOT NULL,
    quantity           INTEGER      NOT NULL,
    nhan_vien_thuc_hien VARCHAR(255) NOT NULL,
    note               TEXT,
    created_by         VARCHAR(50)  NOT NULL,
    updated_by         VARCHAR(50)  NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
);