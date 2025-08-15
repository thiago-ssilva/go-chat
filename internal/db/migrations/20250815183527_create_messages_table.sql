-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) UNIQUE NOT NULL,
	message TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
