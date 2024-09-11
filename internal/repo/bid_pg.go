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
	const query = `INSERT INTO bid (name, description, status, tender_id, organization_id, creator_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Bid](rows)
}

func (r *bidPG) GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {
	const query = `SELECT * FROM bid WHERE id = $1 ORDER BY version DESC`

	rows, err := r.Pool.Query(ctx, query, bidID)
	if err != nil {
		return nil, err
	}

	return collectOneRow[entity.Bid](rows)
}

func (r *bidPG) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.Bid, error) {
	const query = `SELECT * FROM
		(SELECT DISTINCT ON (id) * 
		FROM bid WHERE creator_id = $1 
		ORDER BY id, version DESC
		LIMIT $2 OFFSET $3)
		ORDER BY name ASC`

	rows, err := r.Pool.Query(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Bid])
}

func (r *bidPG) GetByTenderID(ctx context.Context, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error) {
	const query = `SELECT * FROM
		(SELECT DISTINCT ON (id) * 
		FROM bid WHERE tender_id = $1 
		ORDER BY id, version DESC
		LIMIT $2 OFFSET $3)
		ORDER BY name ASC`

	rows, err := r.Pool.Query(ctx, query, tenderID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Bid])
}

func (r *bidPG) Update(ctx context.Context, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM bid WHERE id = $1 ORDER BY version DESC`
	rows, err := tx.Query(ctx, selectQuery, bidID)
	if err != nil {
		return nil, err
	}

	bid, err := collectOneRow[entity.Bid](rows)
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

	const insertQuery = `INSERT INTO bid (id, name, description, status, tender_id, organization_id, creator_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID, bid.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Bid](rows)
}

func (r *bidPG) UpdateStatus(ctx context.Context, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error) {
	const query = `UPDATE bid 
	SET status = $2 
	WHERE id = $1 AND version = (SELECT MAX(version) FROM bid WHERE id = $1) 
	RETURNING *`

	rows, err := r.Pool.Query(ctx, query, bidID, status)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Bid](rows)
}

func (r *bidPG) Rollback(ctx context.Context, bidID uuid.UUID, version int) (*entity.Bid, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM bid WHERE id = $1 AND version = $2`
	rows, err := tx.Query(ctx, selectQuery, bidID, version)
	if err != nil {
		return nil, err
	}

	bid, err := collectExactlyOneRow[entity.Bid](rows)
	if err != nil {
		return nil, err
	}

	bid.Version++

	const insertQuery = `INSERT INTO bid (id, name, description, status, tender_id, organization_id, creator_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID, bid.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Bid](rows)
}
