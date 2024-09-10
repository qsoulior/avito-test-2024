package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"github.com/google/uuid"
)

type employeeV1 struct {
	employeeRepo repo.Employee
}

func NewEmployeeV1(employee repo.Employee) Employee {
	if employee == nil {
		return nil
	}
	return &employeeV1{employee}
}

func (s *employeeV1) GetUser(ctx context.Context, username string) (*entity.Employee, error) {
	employee, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrEmployeeUnauthorized
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return employee, nil
}

func (s *employeeV1) GetEmployee(ctx context.Context, username string, organizationID uuid.UUID) (*entity.Employee, error) {
	employee, err := s.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	err = s.employeeRepo.HasOrganization(ctx, employee.ID, organizationID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrEmployeeForbidden
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return employee, nil
}
