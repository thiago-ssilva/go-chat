-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_username_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages ADD CONSTRAINT messages_username_key UNIQUE (username);
-- +goose StatementEnd
