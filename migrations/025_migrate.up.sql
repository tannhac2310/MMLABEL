drop table customers;

create table customers
(
    id VARCHAR(50) not null,
    name VARCHAR(255) not null,
    tax VARCHAR(255) null,
    code VARCHAR(255) not null,
    country VARCHAR(255) not null,
    province VARCHAR(255) not null,
    address VARCHAR(255) not null,
    phone_number VARCHAR(50) not null,
    fax VARCHAR(255) null,
    company_website VARCHAR(255) null,
    company_phone VARCHAR(50) null,
    company_email VARCHAR(50) null,
    contact_person_name VARCHAR(255) not null,
    contact_person_email VARCHAR(255) not null,
    contact_person_phone VARCHAR(50) not null,
    contact_person_role VARCHAR(255) not null,
    note VARCHAR(255) null,
    status smallint default 1,
    created_by VARCHAR(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz,
    constraint pk_customers primary key (id asc)
);

alter table production_plans add column product_name varchar(500) not null default '';
alter table production_plans add column product_code varchar(255) not null default '';
alter table production_plans add qty_paper int8 null default 0;
alter table production_plans add qty_finished int8 null default 0;
alter table production_plans add qty_delivered int8 null default 0;
