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
    x INT NOT NULL,
    y INT NOT NULL,
    type TEXT NOT NULL,
    parent_id TEXT REFERENCES nodes (id) ON DELETE CASCADE,
    assginer_id INT,
    title TEXT,
    data JSONB,
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
    title TEXT NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    w INT NOT NULL,
    h INT NOT NULL,
    parent_id TEXT,
    type TEXT,
    workflow_id INT NOT NULL REFERENCES workflows (id)
);

-- +goose Down
DROP TABLE node_groups;

DROP TABLE node_connections;

DROP TABLE nodes;

DROP TABLE workflows;
