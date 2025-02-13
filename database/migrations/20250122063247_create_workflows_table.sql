-- +goose Up
CREATE TABLE workflows (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    category_id INT NOT NULL,
    version INT NOT NULL,
    description TEXT NOT NULL,
    decoration TEXT NOT NULL,
    form_id INT REFERENCES forms (id)
);

CREATE TABLE nodes (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    x NUMERIC NOT NULL,
    y NUMERIC NOT NULL,
    width NUMERIC NOT NULL,
    height NUMERIC NOT NULL,
    type TEXT NOT NULL,
    parent_id TEXT,
    summary TEXT,
    end_type TEXT,
    key TEXT,
    workflow_id INT NOT NULL REFERENCES workflows (id)
);

CREATE TABLE node_connections (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    type TEXT NOT NULL,
    from_node_id TEXT NOT NULL REFERENCES nodes (id),
    to_node_id TEXT NOT NULL REFERENCES nodes (id),
    workflow_id INT NOT NULL REFERENCES workflows (id)
);

CREATE TABLE node_groups (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,
    summary TEXT NOT NULL,
    x NUMERIC NOT NULL,
    y NUMERIC NOT NULL,
    width NUMERIC NOT NULL,
    height NUMERIC NOT NULL,
    ticket_id TEXT,
    key TEXT,
    type TEXT,
    workflow_id INT NOT NULL REFERENCES workflows (id)
);

-- +goose Down
DROP TABLE node_groups;

DROP TABLE node_connections;

DROP TABLE nodes;

DROP TABLE workflows;
