-- +goose Up
-- +goose StatementBegin
ALTER TABLE entry
ADD COLUMN key text;  
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE entry
DROP COLUMN key; 
-- +goose StatementEnd
