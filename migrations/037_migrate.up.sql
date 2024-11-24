CREATE TABLE orders
(
    id                    VARCHAR(50) PRIMARY KEY,
    title                 VARCHAR(255)  NOT NULL,
    code                  VARCHAR(50)   NOT NULL,
    sale_name             VARCHAR(255)  NOT NULL,
    sale_admin_name       VARCHAR(255)  NOT NULL,
    product_code          VARCHAR(50)   NOT NULL,
    product_name          VARCHAR(1000) NOT NULL,
    customer_id           VARCHAR(50)   NOT NULL,
    customer_product_code VARCHAR(50)   NOT NULL,
    customer_product_name VARCHAR(255)  NOT NULL,
    status                VARCHAR(50)   NOT NULL,
    created_by            VARCHAR(50)   NOT NULL,
    updated_by            VARCHAR(50)   NOT NULL,
    created_at            TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
create table order_items
(
    id                         VARCHAR(50) PRIMARY KEY,
    production_plan_product_id VARCHAR(50) NOT NULL,
    production_plan_id         VARCHAR(50) NOT NULL,
    production_quantity        INTEGER,                                                           -- Số lượng sản xuất -- Số đơn đặt hàng/Đơn hàng? = Số lượng sản xuất/Số lượng (1)
    quantity                   INTEGER,                                                           -- Số lượng (1)
    unit_price                 NUMERIC(10, 2),                                                    -- Đơn giá (2)
    total_amount               NUMERIC(10, 2) GENERATED ALWAYS AS (quantity * unit_price) STORED, -- Thành tiền
    delivered_quantity         INT         NOT NULL DEFAULT 0,
    estimated_delivery_date    TIMESTAMP,
    delivered_date             TIMESTAMP,
    status                     VARCHAR(50) NOT NULL,
    attachment                 jsonb,
    note                       TEXT,
    created_by                 VARCHAR(50) NOT NULL,
    updated_by                 VARCHAR(50) NOT NULL,
    created_at                 TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);