package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type tenderPG struct {
	*postgres.Postgres
}

func NewTenderPG(pg *postgres.Postgres) Tender {
	if pg == nil {
		return nil
	}
	return &tenderPG{pg}
}

func (r *tenderPG) Create(ctx context.Context, tender entity.Tender) (*entity.Tender, error) {
	const query = `INSERT INTO tender (name, description, service_type, status, organization_id, creator_username) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorUsername)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Tender])
	if err != nil {
		return nil, err
	}

	return row, err
}

func (r *tenderPG) GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	const query = `SELECT * FROM tender WHERE id = $1`

	rows, err := r.Pool.Query(ctx, query, tenderID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Tender])
}

func (r *tenderPG) GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error) {
	const query = `SELECT * FROM tender WHERE creator_username = $1 LIMIT $2 OFFSET $3`
	rows, err := r.Pool.Query(ctx, query, username, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Tender])
}

func (r *tenderPG) Update(ctx context.Context, tenderID uuid.UUID, username string, data TenderData) (*entity.Tender, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM tender WHERE id = $1 AND creator_username = $2`
	rows, err := tx.Query(ctx, selectQuery, tenderID, username)
	if err != nil {
		return nil, err
	}

	tender, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Tender])
	if err != nil {
		return nil, err
	}

	if data.Name != nil {
		tender.Name = *data.Name
	}

	if data.Description != nil {
		tender.Description = *data.Description
	}

	if data.ServiceType != nil {
		tender.ServiceType = *data.ServiceType
	}

	const updateQuery = `UPDATE tender 
		SET name = $3, description = $4, service_type = $5 
		WHERE id = $1 AND creator_username = $2 RETURNING *`

	rows, err = tx.Query(ctx, updateQuery, tender.ID, tender.CreatorUsername, tender.Name, tender.Description, tender.ServiceType)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Tender])
}

func (r *tenderPG) UpdateStatus(ctx context.Context, tenderID uuid.UUID, username string, status entity.TenderStatus) (*entity.Tender, error) {
	const updateQuery = `UPDATE tender 
		SET status = $3 
		WHERE id = $1 AND creator_username = $2 RETURNING *`

	rows, err := r.Pool.Query(ctx, updateQuery, tenderID, username, status)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Tender])
}
