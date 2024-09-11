package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool     *pgxpool.Pool
	maxConns int32
}

func New(ctx context.Context, uri string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{}

	for _, opt := range opts {
		opt(pg)
	}

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	if pg.maxConns >= 1 {
		cfg.MaxConns = pg.maxConns
	}

	pg.Pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pg.Pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
