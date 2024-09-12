package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type organizationPG struct {
	*postgres.Postgres
}

func NewOrganizationPG(pg *postgres.Postgres) Organization {
	if pg == nil {
		return nil
	}
	return &organizationPG{pg}
}

func (r *organizationPG) GetByID(ctx context.Context, organizationID uuid.UUID) (*entity.Organization, error) {
	const query = `SELECT * FROM organization WHERE id = $1`

	rows, err := r.Pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Organization](rows)
}

func (r *organizationPG) GetByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]entity.Organization, error) {
	const query = `SELECT o.id, o.name, o.description, o.type, o.created_at, o.updated_at
		FROM organization o JOIN organization_responsible r ON o.id = r.organization_id
		WHERE r.user_id = $1`

	rows, err := r.Pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Organization])
}
