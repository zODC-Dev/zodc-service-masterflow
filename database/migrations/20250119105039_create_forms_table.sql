-- +goose Up
CREATE TABLE forms (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    file_name TEXT NOT NULL,
    title TEXT NOT NULL,
    category_id INT,
    version INT NOT NULL,
    template_id INT,
    data_sheet JSONB,
    description TEXT NOT NULL,
    decoration TEXT NOT NULL,
);

CREATE TABLE form_fields (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    field_id TEXT NOT NULL,
    icon TEXT NOT NULL,
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    field_name TEXT NOT NULL,
    field_type TEXT NOT NULL,
    required BOOLEAN NOT NULL,
    advanced_options JSONB NOT NULL,
    col_num INT NOT NULL,
    form_id INT NOT NULL REFERENCES forms (id)
);

-- +goose Down
DROP TABLE form_fields;

DROP TABLE forms;
