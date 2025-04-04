-- +goose Up
CREATE TABLE workflows (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    user_id INT NOT NULL,

    -- Info
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    category_id INT NOT NULL,
    description TEXT NOT NULL,
    decoration TEXT NOT NULL,

    project_key TEXT,

    current_version INT NOT NULL DEFAULT 1,

    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE workflow_versions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

     -- Version
    version INT NOT NULL,
    has_sub_workflow BOOLEAN NOT NULL,

    workflow_id INT NOT NULL REFERENCES workflows (id) ON DELETE CASCADE
);

CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    user_id INT NOT NULL,

    key SERIAL NOT NULL,

    last_update_user_id INT NOT NULL,

    status TEXT NOT NULL, -- TO_DO, IN_PROCESS , COMPLETED, CANCELED, TERMINATED
    title TEXT NOT NULL,

    is_template BOOLEAN NOT NULL DEFAULT FALSE,

    sprint_id INT,

    parent_id INT,

    progress REAL NOT NULL DEFAULT 0.0,

    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    canceled_at TIMESTAMP,
    terminated_at TIMESTAMP,


    -- Foreign Key
    workflow_version_id INT NOT NULL REFERENCES workflow_versions (id) ON DELETE CASCADE
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

    key SERIAL NOT NULL,
    jira_key TEXT,

    -- Info
    title TEXT NOT NULL,

    assignee_id INT,

    sub_request_id INT REFERENCES requests (id) ON DELETE CASCADE,    
    
    type TEXT NOT NULL, -- start, end, bug, task, approve, sub_workflow, story, input, noti, group, condition
    status TEXT NOT NULL, -- TO_DO, IN_PROCESSING, COMPLETED, OVERDUE

    is_current BOOLEAN NOT NULL DEFAULT false,

    -- estimate
    estimate_point INT,
    planned_start_time TIMESTAMP,
    planned_end_time TIMESTAMP,
    actual_start_time TIMESTAMP,
    actual_end_time TIMESTAMP,

    -- notification node
    body TEXT,
    subject TEXT,

    -- approve node
    is_approved BOOLEAN NOT NULL DEFAULT false,

    -- end node
    end_type TEXT,

    -- task node
    task_assigned_requester BOOLEAN NOT NULL DEFAULT false,
    task_assigned_assignee BOOLEAN NOT NULL DEFAULT false,
    task_assigned_participants BOOLEAN NOT NULL DEFAULT false,

    task_started_requester BOOLEAN NOT NULL DEFAULT false,  
    task_started_assignee BOOLEAN NOT NULL DEFAULT false,
    task_started_participants BOOLEAN NOT NULL DEFAULT false,

    task_completed_requester BOOLEAN NOT NULL DEFAULT false,
    task_completed_assignee BOOLEAN NOT NULL DEFAULT false,
    task_completed_participants BOOLEAN NOT NULL DEFAULT false,

    parent_id TEXT REFERENCES nodes (id) ON DELETE CASCADE,

    -- Foreign Key
    request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE,
    form_template_id INT REFERENCES form_templates (id) ON DELETE CASCADE,
    form_data_id INT REFERENCES form_data (id) ON DELETE CASCADE
);

CREATE TABLE node_condition_destinations (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    destination_node_id TEXT NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,

    is_true BOOLEAN NOT NULL,

    node_id TEXT NOT NULL REFERENCES nodes (id) ON DELETE CASCADE
);

CREATE TABLE node_forms (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    permission TEXT NOT NULL,
    key TEXT NOT NULL, -- FE GENERATE UUID

    -- For FE
    option_key TEXT,
    from_user_id INT,
    from_form_attached_position INT,
    is_original BOOLEAN NOT NULL DEFAULT false,

    -- Form
    data_id TEXT NOT NULL,
    template_id INT NOT NULL REFERENCES form_templates (id) ON DELETE CASCADE,

    node_id TEXT NOT NULL REFERENCES nodes (id) ON DELETE CASCADE
);

CREATE TABLE node_form_approve_users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    user_id INT NOT NULL,

    node_form_id INT NOT NULL REFERENCES node_forms (id) ON DELETE CASCADE
);

CREATE TABLE connections (
    id TEXT PRIMARY KEY,

    created_at TIMESTAMP DEFAULT now () NOT NULL,
    updated_at TIMESTAMP DEFAULT now () NOT NULL,
    deleted_at TIMESTAMP,

    from_node_id TEXT NOT NULL REFERENCES nodes (id),
    to_node_id TEXT NOT NULL REFERENCES nodes (id),

    text TEXT,

    is_completed BOOLEAN NOT NULL DEFAULT false, 

    -- Foreign Key
    request_id INT NOT NULL REFERENCES requests (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE node_form_approve_users;

DROP TABLE node_forms;

DROP TABLE node_condition_destinations;

DROP TABLE connections;

DROP TABLE nodes;

DROP TABLE requests;

DROP TABLE workflow_versions;

DROP TABLE workflows;
