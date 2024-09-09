package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type employeePG struct {
	*postgres.Postgres
}

func NewEmployeePG(pg *postgres.Postgres) Employee {
	if pg == nil {
		return nil
	}
	return &employeePG{pg}
}

func (r *employeePG) GetByID(ctx context.Context, employeeID uuid.UUID) (*entity.Employee, error) {
	const query = `SELECT * FROM employee WHERE id = $1`

	rows, err := r.Pool.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Employee])
}

func (r *employeePG) GetByUsername(ctx context.Context, username string) (*entity.Employee, error) {
	const query = `SELECT * FROM employee WHERE username = $1`

	rows, err := r.Pool.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Employee])
}

func (r *employeePG) GetByUsernameAndOrganizationID(ctx context.Context, username string, organizationID uuid.UUID) (*entity.Employee, error) {
	const query = `SELECT (e.id, e.username, e.first_name, e.last_name, e.created_at, e.updated_at) 
		FROM employee e JOIN organization_responsible r ON e.id = r.user_id
		WHERE e.username = $1 AND r.organization_id = $2`

	rows, err := r.Pool.Query(ctx, query, username, organizationID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Employee])
}
