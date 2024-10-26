-- create table production_order_stage_responsible

CREATE TABLE master_data
(
    id          VARCHAR(50) NOT NULL,
    type        string      NOT NULL,
    name        string      NOT NULL,
    description string,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    created_by  VARCHAR(50),
    updated_by  VARCHAR(50),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT pk_master_data PRIMARY KEY (id ASC)
);

create index idx_master_data_type on master_data (type);

CREATE TABLE master_data_user_field
(
    id             VARCHAR(50) NOT NULL,
    master_data_id VARCHAR(50) NOT NULL,
    field_name     string      NOT NULL,
    field_value    string      NOT NULL,
    data           jsonb,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    created_by     VARCHAR(50),
    updated_by     VARCHAR(50),
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT pk_master_data_user_field PRIMARY KEY (id ASC)
);

create index idx_master_data_user_field_master_data_id on master_data_user_field (master_data_id);

CREATE TABLE products
( -- sản phẩm nguồn
    id          VARCHAR(50) NOT NULL,
    name        string      NOT NULL,
    code        string      NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    sale_id    VARCHAR(50),
    description string,
    data        jsonb,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    created_by  VARCHAR(50),
    updated_by  VARCHAR(50),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT pk_products PRIMARY KEY (id ASC)
);

CREATE TABLE production_plan_products
( -- sản phẩm kế hoạch sản xuất
    id                       VARCHAR(50) NOT NULL,
    production_plan_id VARCHAR(50) NOT NULL,
    product_id               VARCHAR(50) NOT NULL,
    quantity                 INT         NOT NULL,
    created_at               TIMESTAMPTZ NOT NULL,
    updated_at               TIMESTAMPTZ NOT NULL,
    created_by               VARCHAR(50) NOT NULL,
    updated_by               VARCHAR(50) NOT NULL,
    CONSTRAINT pk_production_plan_products PRIMARY KEY (id ASC)
);


ALTER TABLE production_plans drop column customer_id;
ALTER TABLE production_plans drop column sales_id;