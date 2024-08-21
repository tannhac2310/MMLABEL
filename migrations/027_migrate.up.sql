create table workflow_templates
(
    id         VARCHAR(50)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    config     JSONB,
    status     SMALLINT              DEFAULT 1,
    created_by VARCHAR(50)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()::TIMESTAMPTZ,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()::TIMESTAMPTZ,
    updated_by VARCHAR(50)  NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT pk_workflow_stage_templates PRIMARY KEY (id ASC)
);