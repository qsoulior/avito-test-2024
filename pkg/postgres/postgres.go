package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool      *pgxpool.Pool
	maxConns  int32
	dataTypes []string
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

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, typeName := range pg.dataTypes {
			dataType, err := conn.LoadType(ctx, typeName)
			if err != nil {
				return err
			}
			conn.TypeMap().RegisterType(dataType)

			dataType, err = conn.LoadType(ctx, fmt.Sprintf("_%s", typeName))
			if err != nil {
				return err
			}
			conn.TypeMap().RegisterType(dataType)
		}
		return nil
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
