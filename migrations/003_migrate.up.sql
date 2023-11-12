CREATE TABLE customers
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    avatar       VARCHAR(512),
    phone_number VARCHAR(50),
    email        VARCHAR(255),
    status       SMALLINT              DEFAULT 1,
    type         SMALLINT              DEFAULT 1,
    address      VARCHAR(512),
    created_by   VARCHAR(50)  NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT pk_customers PRIMARY KEY (id ASC)
);
--
-- CREATE TABLE employees
-- (
--     id           VARCHAR(50)  NOT NULL,
--     name         VARCHAR(255) NOT NULL,
--     avatar       VARCHAR(512),
--     phone_number VARCHAR(50),
--     email        VARCHAR(255),
--     status       SMALLINT              DEFAULT 1,
--     type         SMALLINT              DEFAULT 1,
--     address      VARCHAR(512),
--     created_by   VARCHAR(50)  NOT NULL,
--     created_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
--     updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
--     deleted_at   TIMESTAMPTZ,
--     CONSTRAINT pk_employees PRIMARY KEY (id ASC)
-- );

-- SỐ LỆNH SX	KINH DOANH	KHÁCH HÀNG	TÊN /MÃ SP	TÊN NVL 	SỐ MÉT	SỐ LƯỢNG TỜ	SỐ LƯỢNG TP	SỐ LƯỢNG GIAO	NGÀY BẮT ĐẦU-KẾT THÚC	SỰ CỐ	GHI CHÚ	HÌNH ẢNH SẢN PHẨM
CREATE TABLE production_orders
(
    id                      VARCHAR(50)  NOT NULL,
    product_code            VARCHAR(255) NOT NULL,
    product_name            VARCHAR(255) NOT NULL,
    customer_id             VARCHAR(50)  NOT NULL,
    sales_id                VARCHAR(50)  NOT NULL,
    qty_paper               INT                   DEFAULT NULL,
    qty_finished            INT                   DEFAULT NULL,
    qty_delivered           INT                   DEFAULT NULL,
    planned_production_date TIMESTAMPTZ  NOT NULL,
    delivery_date           TIMESTAMPTZ  NOT NULL,
    delivery_image          VARCHAR(255) NULL,
    status                  SMALLINT              DEFAULT 1,
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
    option_id  VARCHAR(50)  NOT NULL,
    data       JSONB,
    status     SMALLINT              DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_devices PRIMARY KEY (id ASC)
);

CREATE TABLE stages -- công đoạn
(
    id         VARCHAR(50)  NOT NULL,
    parent_id  VARCHAR(50),
    department_code  VARCHAR(50),
    name       VARCHAR(255) NOT NULL,
    short_name VARCHAR(255) NOT NULL,
    code       VARCHAR(255) NOT NULL,
    sorting    INT                   DEFAULT 1,
    error_code TEXT, -- NVL1,NVL2,NVL3
    data       JSONB,
    note       TEXT,
    status     SMALLINT              DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_stages PRIMARY KEY (id ASC)
);


CREATE TABLE production_order_stage
(
    id                  VARCHAR(50)  NOT NULL,
    production_order_id VARCHAR(50)  NOT NULL,              -- lệnh sx
    stage_id            VARCHAR(255) NOT NULL,              -- công đoạn
    started_at          TIMESTAMPTZ,                        -- thời gian bắt đầu start
    completed_at        TIMESTAMPTZ,                        -- thời gian kết thúc
    status              SMALLINT              DEFAULT 1,    -- Tạm dừng/Đang chạy
    condition           VARCHAR(10)           DEFAULT NULL, -- Chờ PC/Chờ SX/Đang Sản xuất/Chuyển PO/Hoàn thành SX/SS Vận chuyển
    note                TEXT,
    data                JSONB,                              -- thông tin bổ sung
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_production_order_stage PRIMARY KEY (id ASC)
);

--- 1 công đoạn + thiết bị -> trạng thái của thiết bị
CREATE TABLE production_order_state_device_assignments
(
    id                        VARCHAR(50)  NOT NULL,
    production_order_stage_id VARCHAR(50)  NOT NULL,
    device_id                 VARCHAR(255) NOT NULL,
    quantity                  INT                   DEFAULT 0,
    process_status            SMALLINT, -- null not set, 1: success, 0: failed
    status                    SMALLINT              DEFAULT 1,
    settings                  JSONB,    -- thông tin cấu hình máy lưu bằng jsonb
    note                      TEXT,
    created_at                TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at                TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at                TIMESTAMPTZ,
    CONSTRAINT pk_production_order_state_device_assignments PRIMARY KEY (id ASC),
    CONSTRAINT unique_production_order_state_assignments UNIQUE (production_order_stage_id, device_id)
);

--- phân công của 1 lệnh sx tại công đoạn X
CREATE TABLE production_order_state_device_employee_assignments
(
    id                       VARCHAR(50) NOT NULL,
    pos_device_assignment_id VARCHAR(50),
    employee_id              VARCHAR(50),
    note                     TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at               TIMESTAMPTZ,
    CONSTRAINT pk_production_order_state_device_employee_assignments PRIMARY KEY (id ASC),
    CONSTRAINT unique_production_order_state_assignments UNIQUE (pos_device_assignment_id, employee_id)
);

--- danh muc
CREATE TABLE options
(
    id         VARCHAR(50)  NOT NULL,
    entity     VARCHAR(255) NOT NULL, -- loại DanhMucChung: MaDM,TenDM,LoaiDM(MaLoi,TinhTrang,CongDoan, DonVi)
    code       VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    data       JSONB,
    status     SMALLINT              DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_options PRIMARY KEY (id ASC)
);
-- exec /cockroach/cockroach start --join=${STATEFULSET_NAME}-0.${STATEFULSET_FQDN}:26257 --advertise-host=$(hostname).${STATEFULSET_FQDN} --certs-dir=/cockroach/cockroach-certs/ --http-port=8080 --port=26257 --cach │
-- │ e=1Gi --max-sql-memory=1Gi --logtostderr=INFO