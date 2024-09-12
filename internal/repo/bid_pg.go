package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// bidPG.
type bidPG struct {
	*postgres.Postgres
}

func NewBidPG(pg *postgres.Postgres) Bid {
	if pg == nil {
		return nil
	}
	return &bidPG{pg}
}

func (r *bidPG) HasByCreatorID(ctx context.Context, creatorID uuid.UUID, tenderID uuid.UUID) error {
	const query = `SELECT * FROM bid WHERE creator_id = $1 AND tender_id = $2`

	rows, err := r.Pool.Query(ctx, query, creatorID, tenderID)
	if err != nil {
		return err
	}

	if !rows.Next() {
		if rows.Err() == nil {
			return ErrNoRows
		}
		return rows.Err()
	}

	rows.Close()
	return rows.Err()
}

func (r *bidPG) Create(ctx context.Context, bid entity.Bid) (*entity.Bid, error) {
	const query = `INSERT INTO bid (name, description, status, tender_id, organization_id, creator_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	rows, err := r.Pool.Query(ctx, query,
		bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID)
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
		LIMIT $2 OFFSET $3) AS bid
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
		FROM bid WHERE tender_id = $1 AND status IN ('Published','Approved','Rejected') 
		ORDER BY id, version DESC
		LIMIT $2 OFFSET $3) AS bid
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

	const insertQuery = `INSERT INTO bid 
		(id, name, description, status, tender_id, organization_id, creator_id, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	rows, err = tx.Query(ctx, insertQuery,
		bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID, bid.Version)
	if err != nil {
		return nil, err
	}

	bid, err = collectExactlyOneRow[entity.Bid](rows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return bid, nil
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

	const insertQuery = `INSERT INTO bid
		(id, name, description, status, tender_id, organization_id, creator_id, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT MAX(version) FROM bid WHERE id = $1) + 1) 
		RETURNING *`

	rows, err = tx.Query(ctx, insertQuery,
		bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.CreatorID)
	if err != nil {
		return nil, err
	}

	bid, err = collectExactlyOneRow[entity.Bid](rows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

// bidReviewPG.
type bidReviewPG struct {
	*postgres.Postgres
}

func NewBidReviewPG(pg *postgres.Postgres) BidReview {
	if pg == nil {
		return nil
	}
	return &bidReviewPG{pg}
}

func (r *bidReviewPG) Create(ctx context.Context, review entity.BidReview) (*entity.BidReview, error) {
	const query = `INSERT INTO bid_review (description, bid_id, organization_id, creator_id) 
		VALUES ($1, $2, $3, $4) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, review.Description, review.BidID, review.OrganizationID, review.CreatorID)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.BidReview](rows)
}

func (r *bidReviewPG) GetByBidCreatorID(ctx context.Context,
	creatorID uuid.UUID, limit int, offset int) ([]entity.BidReview, error) {
	const query = `SELECT bid_review.id, description, bid_id, organization_id, creator_id, created_at FROM 
		(SELECT DISTINCT ON (id) id 
		FROM bid WHERE creator_id = $1 
		ORDER BY id, version DESC) as bid
		JOIN bid_review ON bid.id = bid_review.bid_id
		LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.BidReview])
}

// bidDecisionPG.
type bidDecisionPG struct {
	*postgres.Postgres
}

func NewBidDecisionPG(pg *postgres.Postgres) BidDecision {
	if pg == nil {
		return nil
	}
	return &bidDecisionPG{pg}
}

func (r *bidDecisionPG) Create(ctx context.Context, decision entity.BidDecision) (*entity.BidDecision, error) {
	const query = `INSERT INTO bid_decision (bid_id, type, organization_id, creator_id) 
		VALUES ($1, $2, $3, $4) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, decision.BidID, decision.Type, decision.OrganizationID, decision.CreatorID)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.BidDecision](rows)
}

func (r *bidDecisionPG) GetByBidID(ctx context.Context,
	bidID uuid.UUID, organizationID uuid.UUID, decisionType *entity.BidStatus) ([]entity.BidDecision, error) {
	const query = `SELECT DISTINCT ON (creator_id) * 
		FROM bid_decision WHERE bid_id = $1 AND organization_id = $2 AND ($3::bid_decision_type IS NULL OR type = $3)
		ORDER BY creator_id, created_at DESC`

	rows, err := r.Pool.Query(ctx, query, bidID, organizationID, decisionType)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.BidDecision])
}
