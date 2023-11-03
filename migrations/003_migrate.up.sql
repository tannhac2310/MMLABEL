CREATE TABLE customers
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    status       SMALLINT DEFAULT 1,
    type         SMALLINT DEFAULT 1,
    address      VARCHAR(512),
    created_by   VARCHAR(50)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_customers PRIMARY KEY (id ASC)
);

CREATE TABLE employees
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    status       SMALLINT    DEFAULT 1,
    type         SMALLINT    DEFAULT 1,
    address      VARCHAR(512),
    created_by   VARCHAR(50)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_employees PRIMARY KEY (id ASC)
);

CREATE TABLE production_orders
(
    id                      VARCHAR(50)  NOT NULL,
    product_code            VARCHAR(255) NOT NULL,
    product_name            VARCHAR(255) NOT NULL,
    customer_id             VARCHAR(50)  NOT NULL REFERENCES customers(id),
    planned_production_date TIMESTAMPTZ  NOT NULL,
    delivery_date           TIMESTAMPTZ  NOT NULL,
    responsible             STRING[], -- INSERT INTO a VALUES (ARRAY['sky', 'road', 'car']);
    status                  SMALLINT     DEFAULT 1,
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
    option_id VARCHAR(50) NOT NULL REFERENCES options(id),
    data       JSONB,
    status     SMALLINT     DEFAULT 1,
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
    status     SMALLINT     DEFAULT 1,
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
    started_at          TIMESTAMPTZ,           -- thời gian bắt đầu start
    completed_at        TIMESTAMPTZ,           -- thời gian kết thúc
    status              SMALLINT    DEFAULT 1, -- Tạm dừng/Đang chạy
    condition           VARCHAR(10)           DEFAULT NULL, -- Chờ PC/Chờ SX/Đang Sản xuất/Chuyển PO/Hoàn thành SX/SS Vận chuyển
    note                TEXT,
    data                JSONB,                 -- thông tin bổ sung
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_production_order_stage PRIMARY KEY (id ASC)
);

--- phân công của 1 lệnh sx tại công đoạn X
CREATE TABLE production_order_state_device_assignments
(
    id                          VARCHAR(50)  NOT NULL,
    production_order_stage_id   VARCHAR(50)  NOT NULL,
    device_id                   VARCHAR(255) NOT NULL REFERENCES devices(id),
    quantity                    INT DEFAULT 0,
    sucessful                   INT DEFAULT 0,
    unsucessful                 INT DEFAULT 0,
    status                      SMALLINT DEFAULT 1,
    settings                    JSONB,     -- thông tin cấu hình máy lưu bằng jsonb
    note                        TEXT,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at                  TIMESTAMPTZ,
    CONSTRAINT pk_production_order_state_device_assignments PRIMARY KEY (id ASC),
    CONSTRAINT unique_production_order_state_assignments UNIQUE (production_order_stage_id, device_id)
);

--- phân công của 1 lệnh sx tại công đoạn X
CREATE TABLE production_order_state_device_employee_assignments
(
    id                          VARCHAR(50)  NOT NULL,
    pos_device_assignment_id    VARCHAR(50) REFERENCES production_order_state_device_assignments(id),
    employee_id                 VARCHAR(50),
    note                        TEXT,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at                  TIMESTAMPTZ,
    CONSTRAINT pk_production_order_state_device_employee_assignments PRIMARY KEY (id ASC),
    CONSTRAINT unique_production_order_state_assignments UNIQUE (pos_device_assignment_id,employee_id)
);

--- danh muc
CREATE TABLE options
(
    id         VARCHAR(50)  NOT NULL,
    entity     VARCHAR(255) NOT NULL, -- loại DanhMucChung: MaDM,TenDM,LoaiDM(MaLoi,TinhTrang,CongDoan, DonVi)
    code       VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     SMALLINT     DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_options PRIMARY KEY (id ASC)
);
