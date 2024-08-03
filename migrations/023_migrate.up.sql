CREATE TABLE ink_mixing
(
    id               VARCHAR(50)  NOT NULL,
    name             VARCHAR(255) NOT NULL,
    code             VARCHAR(255) NOT NULL,
    ink_id           VARCHAR(50)  NOT NULL,
    mixing_date      TIMESTAMPTZ,
    description      TEXT,
    status           SMALLINT              DEFAULT 1,
    data             JSONB,
    created_by       VARCHAR(50)  NOT NULL,
    updated_by       VARCHAR(50)  NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now()::TIMESTAMPTZ,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()::TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT ink_mixing_pk PRIMARY KEY (id)
);

CREATE TABLE ink_mixing_detail
(
    id            VARCHAR(50) NOT NULL,
    ink_mixing_id VARCHAR(50) NOT NULL,
    ink_id        VARCHAR(50) NOT NULL,
    quantity      FLOAT       NOT NULL,
    description   TEXT,
    data          JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT ink_mixing_detail_pk PRIMARY KEY (id)
);

ALTER TABLE ink ADD COLUMN mixing_id VARCHAR(50);
