-- +goose Up
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP,

    user_id INT NOT NULL,
    content TEXT NOT NULL,

    node_id TEXT NOT NULL REFERENCES nodes (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comments;
