CREATE TABLE customers
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    status       INT2                  DEFAULT 1,
    type         INT2                  DEFAULT 1,
    address      VARCHAR(512),
    created_by   VARCHAR(50)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_customers PRIMARY KEY (id ASC)
);

CREATE TABLE production_orders
(
    id                      VARCHAR(50)  NOT NULL,
    product_code            VARCHAR(255) NOT NULL,
    product_name            VARCHAR(255) NOT NULL,
    customer_id             VARCHAR(512),
    planned_production_date TIMESTAMPTZ  NOT NULL,
    delivery_date           TIMESTAMPTZ  NOT NULL,
    responsible             STRING[], -- INSERT INTO a VALUES (ARRAY['sky', 'road', 'car']);
    status                  INT2                  DEFAULT 1,
    note                    TEXT,
    created_by              VARCHAR(50)  NOT NULL,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at              TIMESTAMPTZ,
    CONSTRAINT pk_production_orders PRIMARY KEY (id ASC)
);

CREATE TABLE devices
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     INT2                  DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_devices PRIMARY KEY (id ASC)
);

CREATE TABLE stages -- công đoạn
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     INT2                  DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_stages PRIMARY KEY (id ASC)
);


CREATE TABLE production_order_stage
(
    id                  VARCHAR(50)  NOT NULL,
    production_order_id VARCHAR(50)  NOT NULL, -- lệnh sx
    stage_id            VARCHAR(255) NOT NULL, -- công đoạn
    device_id           VARCHAR(255),          --  thiết bị
    started_at          TIMESTAMPTZ,           -- thời gian bắt đầu start
    completed_at        TIMESTAMPTZ,           -- thời gian kết thúc
    status              INT2                  DEFAULT 1,
    note                TEXT,
    data                JSONB,                 -- thông tin bổ sung
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_production_order_stage PRIMARY KEY (id ASC)
);

--- cấu hinhf thiết bị của 1 lệnh sx tại công đoạn X
CREATE TABLE production_order_device_configs
(
    id                  VARCHAR(50)  NOT NULL,
    production_order_id VARCHAR(50)  NOT NULL,
    stage_id            VARCHAR(255) NOT NULL,
    device_id           VARCHAR(255) NOT NULL,
    status              INT2                  DEFAULT 1,
    data                JSONB,                  -- thông tin cấu hình máy lưu bằng jsonb
    note                TEXT,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_production_order_device_configs PRIMARY KEY (id ASC)
);

CREATE TABLE stages -- công đoạn
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    code       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     INT2                  DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_stages PRIMARY KEY (id ASC)
);
--- danh muc
CREATE TABLE options
(
    id         VARCHAR(50)  NOT NULL,
    entity     VARCHAR(255) NOT NULL, -- loại DanhMucChung: MaDM,TenDM,LoaiDM(MaLoi,TinhTrang,CongDoan, DonVi)
    code       VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     INT2                  DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_options PRIMARY KEY (id ASC)
);
