package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type bidPG struct {
	*postgres.Postgres
}

func NewBidPG(pg *postgres.Postgres) Bid {
	if pg == nil {
		return nil
	}
	return &bidPG{pg}
}

func (r *bidPG) Create(ctx context.Context, bid entity.Bid) (*entity.Bid, error) {
	const query = `INSERT INTO bid (name, description, status, tender_id, author_type, author_id, creator_username) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.CreatorUsername)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
}

func (r *bidPG) GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {
	const query = `SELECT * FROM bid WHERE id = $1 ORDER BY version DESC`

	rows, err := r.Pool.Query(ctx, query, bidID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
}

func (r *bidPG) GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error) {
	const query = `SELECT DISTINCT ON (id) * 
		FROM bid WHERE creator_username = $1 
		ORDER_BY version DESC, name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, username, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Bid])
}

func (r *bidPG) GetByTenderID(ctx context.Context, bidID uuid.UUID, limit int, offset int) ([]entity.Bid, error) {
	const query = `SELECT DISTINCT ON (id) * 
		FROM bid WHERE $1 IS NULL OR service_type = $1 
		ORDER_BY version DESC , name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, bidID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Bid])
}

func (r *bidPG) Update(ctx context.Context, bidID uuid.UUID, username string, data BidData) (*entity.Bid, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM bid WHERE id = $1 AND creator_username = $2 ORDER BY version DESC`
	rows, err := tx.Query(ctx, selectQuery, bidID, username)
	if err != nil {
		return nil, err
	}

	bid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
	if err != nil {
		return nil, err
	}

	if data.Name != nil {
		bid.Name = *data.Name
	}

	if data.Description != nil {
		bid.Description = *data.Description
	}

	bid.Version++

	const insertQuery = `INSERT INTO bid (id, name, description, status, tender_id, author_type, author_id, creator_username, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.CreatorUsername, bid.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
}

func (r *bidPG) UpdateStatus(ctx context.Context, bidID uuid.UUID, username string, status entity.BidStatus) (*entity.Bid, error) {
	const query = `UPDATE bid 
	SET status = $3 
	WHERE id = $1 AND creator_username = $2 AND version = (SELECT MAX(version) FROM bid WHERE id = $1 AND creator_username = $2) 
	RETURNING *`

	rows, err := r.Pool.Query(ctx, query, bidID, username, status)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
}

func (r *bidPG) Rollback(ctx context.Context, bidID uuid.UUID, username string, version int) (*entity.Bid, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM bid WHERE id = $1 AND creator_username = $2 AND version = $3`
	rows, err := tx.Query(ctx, selectQuery, bidID, username)
	if err != nil {
		return nil, err
	}

	bid, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
	if err != nil {
		return nil, err
	}

	bid.Version++

	const insertQuery = `INSERT INTO bid (id, name, description, status, tender_id, author_type, author_id, creator_username, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.AuthorType, bid.AuthorID, bid.CreatorUsername, bid.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Bid])
}
