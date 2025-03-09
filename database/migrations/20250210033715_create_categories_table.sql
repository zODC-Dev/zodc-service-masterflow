-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    key TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE workflows
ADD CONSTRAINT workflows_category_id_fkey
FOREIGN KEY (category_id) 
REFERENCES categories (id);

ALTER TABLE form_templates
ADD CONSTRAINT form_templates_category_id_fkey
FOREIGN KEY (category_id) 
REFERENCES categories (id);

-- +goose Down
ALTER TABLE form_templates
DROP CONSTRAINT form_templates_category_id_fkey;

ALTER TABLE workflows
DROP CONSTRAINT workflows_category_id_fkey;

DROP TABLE categories;
