create table production_plans
(
    id VARCHAR(50) NOT NULL,
	customer_id VARCHAR(350) NOT NULL,
	sales_id VARCHAR(50) NOT NULL,
	thumbnail VARCHAR(255) NULL,
	status INT2 NULL DEFAULT 1:::INT8,
	note STRING NULL,
	created_by VARCHAR(50) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
    updated_by VARCHAR(50) NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now():::TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ NULL,
	name VARCHAR(255) NOT NULL,
	CONSTRAINT pk_production_plans PRIMARY KEY (id ASC)
);

create table production_plan_attributes
(
    id VARCHAR(50) not null,
    kind INT2 default 1,
    display_name text not null,
    attribute_value text not null,
    note text,
    data JSONB,
    status smallint default 1,
    created_by VARCHAR(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by VARCHAR(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz,
    constraint pk_production_plan_attributes primary key (id asc)
);

alter table production_plans add column po_stages jsonb default '{}';
alter table production_plans add column stages_info int4 default 1;