-- +goose Up
-- +goose StatementBegin
ALTER TABLE entry ALTER COLUMN book_id SET NOT NULL;
ALTER TABLE entry ALTER COLUMN content SET NOT NULL;
ALTER TABLE entry ALTER COLUMN type SET NOT NULL;

ALTER TABLE book ALTER COLUMN name SET NOT NULL;
-- +goose StatementEnd