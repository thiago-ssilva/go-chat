package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	id        int
	username  string
	content   string
	createdAt time.Time
}

type MessagesRepository struct {
	db *sql.DB
}

func NewMessagesRepository(db *sql.DB) *MessagesRepository {
	return &MessagesRepository{
		db: db,
	}
}

func (r *MessagesRepository) CreateMessage(ctx context.Context, content, username string) (*Message, error) {
	query := `
		INSERT INTO messages (username, content)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	message := &Message{
		username: username,
		content:  content,
	}

	err := r.db.QueryRowContext(ctx, query, message.username, message.content).Scan(
		&message.id,
		&message.createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}

	return message, nil
}

func (r *MessagesRepository) GetAllMessages(ctx context.Context) ([]*Message, error) {
	query := `SELECT id, username, content, created_at FROM messages ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("query all messages: %w", err)
	}

	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var message Message

		err := rows.Scan(
			&message.id,
			&message.username,
			&message.content,
			&message.createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}

	return messages, nil
}
