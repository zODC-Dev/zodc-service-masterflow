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
    has_sub_workflow BOOLEAN NOT NULL,

    status TEXT,

    workflow_id INT NOT NULL REFERENCES workflows (id) ON DELETE CASCADE
);

CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    status TEXT NOT NULL, -- TO_DO, IN_PROCESS , COMPLETED, CANCELED, TERMINATED
    title TEXT NOT NULL,
    description TEXT NOT NULL,

    is_template BOOLEAN NOT NULL DEFAULT FALSE,

    -- Foreign Key
    workflow_version_id INT NOT NULL REFERENCES workflows (id) ON DELETE CASCADE
);

CREATE TABLE nodes (
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

    sub_request_id INT REFERENCES requests (id) ON DELETE CASCADE,    
    
    type TEXT NOT NULL, -- start, end, bug, task, approve, sub_workflow, story, input, noti, group, condition
    status TEXT NOT NULL, -- TO_DO, IN_PROCESSING, COMPLETED, OVERDUE

    -- estimate
    estimate_point INT,
    plan_start_time TIMESTAMP,
    plan_finish_time TIMESTAMP,
    actual_start_time TIMESTAMP,
    actual_finish_time TIMESTAMP,


    parent_id TEXT REFERENCES nodes (id) ON DELETE CASCADE,

    -- Foreign Key
    request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE,

    -- Form
    form_template_id INT REFERENCES form_templates (id) ON DELETE CASCADE,
    form_data_id INT REFERENCES form_data (id) ON DELETE CASCADE
);

CREATE TABLE connections (
    id TEXT PRIMARY KEY,

    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    from_node_id TEXT NOT NULL REFERENCES nodes (id),
    to_node_id TEXT NOT NULL REFERENCES nodes (id),

    type TEXT NOT NULL,

    is_completed BOOLEAN NOT NULL DEFAULT false, 

    -- Foreign Key
    request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE connections;

DROP TABLE nodes;

DROP TABLE requests;

DROP TABLE workflow_versions;

DROP TABLE workflows;
