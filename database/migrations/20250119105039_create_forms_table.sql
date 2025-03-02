-- +goose Up
CREATE TABLE form_templates (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    
    -- Info
    file_name TEXT NOT NULL,
    title TEXT NOT NULL,
    category_id INT,
    data_sheet JSONB,
    description TEXT NOT NULL,
    decoration TEXT NOT NULL,

    tag TEXT NOT NULL DEFAULT 'FORM', -- FORM / TASK, BUG, STORY
    type TEXT NOT NULL DEFAULT 'USER' -- SYSTEM, USER
);

CREATE TABLE form_template_versions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    version INT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT,

    form_template_id INT NOT NULL REFERENCES form_templates (id) ON DELETE CASCADE
);

CREATE TABLE form_template_fields (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    advanced_options JSONB,
    col_num INT NOT NULL,
    required BOOLEAN NOT NULL,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    icon TEXT NOT NULL,
    field_id TEXT NOT NULL,
    field_type TEXT NOT NULL,
    
    field_name TEXT NOT NULL,

    form_template_version_id INT NOT NULL REFERENCES form_template_versions (id) ON DELETE CASCADE
);

CREATE TABLE form_data (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    form_template_version_id INT NOT NULL REFERENCES form_template_versions (id) ON DELETE CASCADE
);

CREATE TABLE form_field_data (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    value TEXT NOT NULL DEFAULT '',

    form_data_id INT NOT NULL REFERENCES form_data (id) ON DELETE CASCADE,
    form_template_field_id INT NOT NULL REFERENCES form_template_fields (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE form_field_data;

DROP TABLE form_data;

DROP TABLE form_template_fields;

DROP TABLE form_template_versions;

DROP TABLE form_templates;
