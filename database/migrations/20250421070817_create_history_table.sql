-- +goose Up
CREATE TABLE histories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP,

    type_action TEXT NOT NULL, -- STATUS, ASSIGNEE, APPROVE, REJECT, START_REQUEST, END_REQUEST, NEW_TASK, CANCEL_REQUEST

    user_id INT,

    request_id INT NOT NULL REFERENCES requests (id),
    node_id TEXT NOT NULL REFERENCES nodes (id),

    from_value TEXT,
    to_value TEXT
);

-- +goose Down
DROP TABLE histories;
