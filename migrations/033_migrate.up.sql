CREATE TABLE master_data_selection
(
    id                  VARCHAR(50) NOT NULL,
    selection_group     VARCHAR(255) NOT NULL DEFAULT ''::VARCHAR,
    display_value       TEXT NOT NULL,
    value               TEXT NOT NULL,
    description         TEXT,
    multiple_choices    INT2,
    sort_order          INT2,
    status              INT2 NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()::TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()::TIMESTAMPTZ,
    created_by          VARCHAR(50),
    updated_by          VARCHAR(50),
    deleted_at          TIMESTAMPTZ,
    CONSTRAINT pk_master_data_selection PRIMARY KEY (id ASC)
);