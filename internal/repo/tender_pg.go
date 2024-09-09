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
	const query = `INSERT INTO tender (name, description, service_type, status, organization_id, creator_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	rows, err := r.Pool.Query(ctx, query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorID)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Tender](rows)
}

func (r *tenderPG) GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	const query = `SELECT * FROM tender WHERE id = $1 ORDER BY version DESC`

	rows, err := r.Pool.Query(ctx, query, tenderID)
	if err != nil {
		return nil, err
	}

	return collectOneRow[entity.Tender](rows)
}

func (r *tenderPG) GetByServiceType(ctx context.Context, serviceType *entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error) {
	const query = `SELECT DISTINCT ON (id) * 
		FROM tender WHERE $1 IS NULL OR service_type = $1 
		ORDER_BY version DESC, name ASC 
		LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, serviceType, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Tender])
}

func (r *tenderPG) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.Tender, error) {
	const query = `SELECT DISTINCT ON (id) * 
		FROM tender WHERE creator_id = $1 
		ORDER_BY version DESC, name ASC 
		LIMIT $2 OFFSET $3`

	rows, err := r.Pool.Query(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Tender])
}

func (r *tenderPG) Update(ctx context.Context, tenderID uuid.UUID, data TenderData) (*entity.Tender, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM tender WHERE id = $1 ORDER BY version DESC`
	rows, err := tx.Query(ctx, selectQuery, tenderID)
	if err != nil {
		return nil, err
	}

	tender, err := collectOneRow[entity.Tender](rows)
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

	tender.Version++

	const insertQuery = `INSERT INTO tender (id, name, description, service_type, status, organization_id, creator_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorID, tender.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Tender](rows)
}

func (r *tenderPG) UpdateStatus(ctx context.Context, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error) {
	const query = `UPDATE tender 
		SET status = $2
		WHERE id = $1 AND version = (SELECT MAX(version) FROM tender WHERE id = $1)
		RETURNING *`

	rows, err := r.Pool.Query(ctx, query, tenderID, status)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Tender](rows)
}

func (r *tenderPG) Rollback(ctx context.Context, tenderID uuid.UUID, version int) (*entity.Tender, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const selectQuery = `SELECT * FROM tender WHERE id = $1 AND version = $2`
	rows, err := tx.Query(ctx, selectQuery, tenderID, version)
	if err != nil {
		return nil, err
	}

	tender, err := collectExactlyOneRow[entity.Tender](rows)
	if err != nil {
		return nil, err
	}

	tender.Version++

	const insertQuery = `INSERT INTO tender (id, name, description, service_type, status, organization_id, creator_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	rows, err = r.Pool.Query(ctx, insertQuery, tender.ID, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorID, tender.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Tender](rows)
}
