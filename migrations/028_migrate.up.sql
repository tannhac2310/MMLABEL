create table comments
(
    id varchar(50) not null,
    user_id varchar(50) not null,
    target_id varchar(50) not null,
    target_type smallint not null default 0,
    content text not null default '',
    created_at timestamptz not null default now()::timestamptz,
    updated_at timestamptz not null default now()::timestamptz,
    deleted_at timestamptz,
    constraint pk_comments primary key (id asc)
);

create table comment_histories (
    id varchar(50) not null,
    comment_id varchar(50) not null,
    content text not null default '',
    created_at timestamptz not null default now()::timestamptz,
    constraint pk_comment_histories primary key (id asc)
);

create table comment_attachments (
    id varchar(50) not null,
    comment_id varchar(50) not null,
    url varchar(500) not null,
    created_at timestamptz not null default now()::timestamptz,
    deleted_at timestamptz,
    constraint pk_comment_attachments primary key (id asc)
);