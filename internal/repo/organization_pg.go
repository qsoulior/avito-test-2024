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

func (r *organizationPG) GetByID(ctx context.Context, organizationID uuid.UUID) (*entity.Organization, error) {
	const query = `SELECT * FROM organization WHERE id = $1`

	rows, err := r.Pool.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[entity.Organization])
}
