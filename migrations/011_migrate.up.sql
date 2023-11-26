CREATE TABLE product_quality
(
    id                  VARCHAR(50) NOT NULL,
    production_order_id VARCHAR(50) NOT NULL,
    product_id          VARCHAR(50),
    defect_type         VARCHAR(255),
    defect_code         VARCHAR(50),
    defect_level        SMALLINT             DEFAULT 1,
    production_stage_id VARCHAR(50),
    defective_quantity  INT, -- Số sản phẩm bị lỗi
    good_quantity       INT, -- Số sản phẩm đạt chất lượng
    description         TEXT,
    created_by          VARCHAR(50) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_product_quality PRIMARY KEY (id ASC)
);