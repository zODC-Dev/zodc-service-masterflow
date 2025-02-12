-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

ALTER TABLE workflows
ADD CONSTRAINT workflows_category_id_fkey
FOREIGN KEY (category_id) 
REFERENCES categories(id);

ALTER TABLE forms
ADD CONSTRAINT forms_category_id_fkey
FOREIGN KEY (category_id) 
REFERENCES categories(id);

-- +goose Down
ALTER TABLE forms
DROP CONSTRAINT forms_category_id_fkey;

ALTER TABLE workflows
DROP CONSTRAINT workflows_category_id_fkey;

DROP TABLE categories;
