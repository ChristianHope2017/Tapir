-- Filename: migrations/000003_create_todo_table.up.sql
CREATE TABLE IF NOT EXISTS todo (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    task text NOT NULL
);