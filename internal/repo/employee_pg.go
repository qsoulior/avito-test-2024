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

	return collectExactlyOneRow[entity.Employee](rows)
}

func (r *employeePG) GetByUsername(ctx context.Context, username string) (*entity.Employee, error) {
	const query = `SELECT * FROM employee WHERE username = $1`

	rows, err := r.Pool.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}

	return collectExactlyOneRow[entity.Employee](rows)
}

func (r *employeePG) GetByOrganization(ctx context.Context, organizationID uuid.UUID) ([]entity.Employee, error) {
	const query = `SELECT (id, username, first_name, last_name, created_at, updated_at)
		 FROM employee JOIN (SELECT (user_id) FROM organization_responsible WHERE organization_id = $1) AS r 
		 ON employee.id = r.user_id`

	rows, err := r.Pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Employee])
}

func (r *employeePG) HasOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) error {
	const query = `SELECT * FROM organization_responsible WHERE user_id = $1 AND organization_id = $2`

	rows, err := r.Pool.Query(ctx, query, userID, organizationID)
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
