create table m_khach_hang (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_thong_so_may_in (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_thong_so_may_khac (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_khung_in (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_phim (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_khuon_be (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_khuon_dap (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);

create table m_nguyen_vat_lieu (
    id varchar(50) not null primary key,
    data jsonb null default '{}',
    created_by varchar(50) not null,
    created_at timestamptz not null default now():::timestamptz,
    updated_by varchar(50) not null,
    updated_at timestamptz not null default now():::timestamptz,
    deleted_at timestamptz
);