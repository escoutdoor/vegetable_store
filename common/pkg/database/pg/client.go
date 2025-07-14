package pg

import (
	"context"
	"fmt"

	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type client struct {
	db database.DB
}

var _ database.Client = (*client)(nil)

func NewClient(ctx context.Context, dsn string) (*client, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	return &client{
		db: NewDB(pool),
	}, nil
}

func (c *client) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *client) DB() database.DB {
	return c.db
}
