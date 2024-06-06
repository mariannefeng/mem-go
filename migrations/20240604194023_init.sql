-- +goose Up
-- +goose StatementBegin
CREATE TABLE book (
    id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE entry (
    id SERIAL PRIMARY KEY,
    book_id INTEGER REFERENCES book (id),
    content TEXT,
    type text
); 
-- +goose StatementEnd
