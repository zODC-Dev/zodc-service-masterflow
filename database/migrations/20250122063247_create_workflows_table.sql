-- +goose Up
CREATE TABLE workflows (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    -- Info
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    category_id INT NOT NULL,
    description TEXT NOT NULL,
    decoration TEXT NOT NULL,

    project_key TEXT
);

CREATE TABLE workflow_versions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

     -- Version
    version INT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT,

    workflow_id INT NOT NULL REFERENCES workflows (id) ON DELETE CASCADE
);

CREATE TABLE workflow_nodes (
    id TEXT PRIMARY KEY,

    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    -- Shape
    x NUMERIC NOT NULL,
    y NUMERIC NOT NULL,
    width NUMERIC NOT NULL,
    height NUMERIC NOT NULL,

    -- Info
    title TEXT,

    assignee_id INT,

    due_in INT,
    end_type TEXT,
    sub_workflow_version_id INT REFERENCES workflow_versions (id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- start, end, bug, task, approve, sub_workflow, story, input, noti, group, condition

    parent_id TEXT REFERENCES workflow_nodes (id) ON DELETE CASCADE,

    -- Foreign Key
    workflow_version_id INT NOT NULL REFERENCES workflow_versions (id) ON DELETE CASCADE,

    -- Form
    form_template_id INT REFERENCES form_templates (id) ON DELETE CASCADE,
    form_data_id INT REFERENCES form_data (id) ON DELETE CASCADE
);

CREATE TABLE workflow_connections (
    id TEXT PRIMARY KEY,

    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    from_workflow_node_id TEXT NOT NULL REFERENCES workflow_nodes (id),
    to_workflow_node_id TEXT NOT NULL REFERENCES workflow_nodes (id),

    type TEXT NOT NULL,

    -- Foreign Key
    workflow_version_id INT NOT NULL REFERENCES workflow_versions (id) ON DELETE CASCADE
);

CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    status TEXT NOT NULL, -- IN_PROCESS , COMPLETED, CANCELED, TERMINATED

    -- Foreign Key
    workflow_version_id INT NOT NULL REFERENCES workflows (id) ON DELETE CASCADE
);

CREATE TABLE request_nodes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    -- COPY OF workflow_nodes


    status TEXT NOT NULL, -- TO_DO, IN_PROCESS, COMPLETED

    request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE,

    --
    form_data_id INT REFERENCES form_data (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE request_nodes;

DROP TABLE requests;

DROP TABLE workflow_connections;

DROP TABLE workflow_nodes;

DROP TABLE workflow_versions;

DROP TABLE workflows;
