package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

var (
	ErrEmployeeUnauthorized = NewTypedError("unauthorized user", ErrorTypeUnauthorized, nil)
	ErrEmployeeForbidden    = NewTypedError("user is not an employee of organization", ErrorTypeForbidden, nil)
)

type Employee interface {
	GetUser(ctx context.Context, username string) (*entity.Employee, error)
	GetEmployee(ctx context.Context, username string, organizationID uuid.UUID) (*entity.Employee, error)
	GetByOrganization(ctx context.Context, organizationID uuid.UUID) ([]entity.Employee, error)
}
