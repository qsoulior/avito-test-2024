package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type Employee interface {
	GetByID(ctx context.Context, employeeID uuid.UUID) (*entity.Employee, error)
	GetByUsername(ctx context.Context, username string) (*entity.Employee, error)
	GetByUsernameAndOrganizationID(ctx context.Context, username string, organizationID uuid.UUID) (*entity.Employee, error)
}
