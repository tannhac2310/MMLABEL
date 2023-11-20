-- create table quản lý khoa mực
-- column: id, name, code, mã sản phẩm áp dụng, Vị trí kệ, vị trí, Nhà sản xuât, số lượng tồn, ngày hết hạn, data  description, created_at, updated_at
CREATE TABLE ink
(
    id              VARCHAR(50)  NOT NULL,
    name            VARCHAR(255) NOT NULL,
    code            VARCHAR(255) NOT NULL,
    product_codes   VARCHAR(255) NOT NULL, -- mã sản phẩm áp dụng
    position        VARCHAR(255) NOT NULL,
    location        VARCHAR(255) NOT NULL,
    manufacturer    VARCHAR(255) NOT NULL,
    color_detail    JSONB,                 -- chi tiết màu sắc
    quantity        FLOAT        NOT NULL,
    expiration_date DATE         NOT NULL,
    description     TEXT,
    data            JSONB,
    status          SMALLINT              DEFAULT 1,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT ink_management_pk PRIMARY KEY (id)
);
-- create table nhập kho mực
CREATE TABLE ink_import
(
    id               VARCHAR(50)  NOT NULL,
    code             VARCHAR(255) NOT NULL,
    import_date      DATE         NOT NULL,
    import_user      VARCHAR(255) NOT NULL,
    import_warehouse VARCHAR(255) NOT NULL,
    export_warehouse VARCHAR(255) NOT NULL,
    description      TEXT,
    status           SMALLINT              DEFAULT 1,
    data             JSONB,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT ink_import_pk PRIMARY KEY (id)
);
-- CREATE TABLE: ink_import_detail
CREATE TABLE ink_import_detail
(
    id            VARCHAR(50) NOT NULL,
    ink_import_id VARCHAR(50) NOT NULL,
    quantity      FLOAT       NOT NULL,
    color_detail  JSONB, -- chi tiết màu sắc
    description   TEXT,
    data          JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT ink_export_detail_pk PRIMARY KEY (id)
);

-- create table xuất kho mực
CREATE TABLE ink_export
(
    id                  VARCHAR(50)  NOT NULL,
    code                VARCHAR(255) NOT NULL,
    production_order_id VARCHAR(50)  NOT NULL, -- xuất từ lệnh sản xuất
    export_date         DATE         NOT NULL,
    export_user         VARCHAR(255) NOT NULL,
    export_warehouse    VARCHAR(255) NOT NULL,
    description         TEXT,
    status              SMALLINT              DEFAULT 1,
    data                JSONB,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT ink_export_pk PRIMARY KEY (id)
);
-- CREATE TABLE: ink_export_detail
CREATE TABLE ink_export_detail
(
    id            VARCHAR(50) NOT NULL,
    ink_export_id VARCHAR(50) NOT NULL,
    quantity      FLOAT       NOT NULL,
    color_detail  JSONB, -- chi tiết màu sắc
    description   TEXT,
    data          JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT ink_export_detail_pk PRIMARY KEY (id)
);

-- create table trả mực về kho
CREATE TABLE ink_return
(
    id               VARCHAR(50)  NOT NULL,
    code             VARCHAR(255) NOT NULL,
    return_date      DATE         NOT NULL,
    return_user      VARCHAR(255) NOT NULL,
    return_warehouse VARCHAR(255) NOT NULL,
    description      TEXT,
    status           SMALLINT              DEFAULT 1,
    data             JSONB,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT ink_return_pk PRIMARY KEY (id)
);
-- CREATE TABLE: ink_return_detail
CREATE TABLE ink_return_detail
(
    id            VARCHAR(50) NOT NULL,
    ink_return_id VARCHAR(50) NOT NULL,
    quantity      FLOAT       NOT NULL,
    color_detail  JSONB, -- chi tiết màu sắc
    description   TEXT,
    data          JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT ink_return_detail_pk PRIMARY KEY (id)
);

-- create table kiểm kho mực
CREATE TABLE ink_inventory
(
    id                  VARCHAR(50)  NOT NULL,
    code                VARCHAR(255) NOT NULL,
    inventory_date      DATE         NOT NULL,
    inventory_user      VARCHAR(255) NOT NULL,
    inventory_warehouse VARCHAR(255) NOT NULL,
    description         TEXT,
    status              SMALLINT              DEFAULT 1,
    data                JSONB,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT ink_inventory_pk PRIMARY KEY (id)
);
-- CREATE TABLE: ink_inventory_detail (chi tiết kiểm kho mực)
CREATE TABLE ink_inventory_detail
(
    id               VARCHAR(50) NOT NULL,
    ink_inventory_id VARCHAR(50) NOT NULL,
    quantity         FLOAT       NOT NULL,
    color_detail     JSONB, -- chi tiết màu sắc
    description      TEXT,
    data             JSONB,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT ink_inventory_detail_pk PRIMARY KEY (id)
);
