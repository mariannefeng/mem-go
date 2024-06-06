-- +goose Up
-- +goose StatementBegin
ALTER TABLE entry ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE book ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
-- +goose StatementEnd