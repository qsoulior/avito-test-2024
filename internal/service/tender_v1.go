package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"github.com/google/uuid"
)

type tenderV1 struct {
	tenderRepo       repo.Tender
	employeeRepo     repo.Employee
	organizationRepo repo.Organization
}

func NewTenderV1(tender repo.Tender, employee repo.Employee, organization repo.Organization) Tender {
	if tender == nil || employee == nil || organization == nil {
		return nil
	}
	return &tenderV1{tender, employee, organization}
}

func (s *tenderV1) getByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	tender, err := s.tenderRepo.GetByID(ctx, tenderID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrTenderNotExist
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return tender, nil
}

func (s *tenderV1) getUser(ctx context.Context, username string) (*entity.Employee, error) {
	employee, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrEmployeeUnauthorized
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return employee, nil
}

func (s *tenderV1) getEmployee(ctx context.Context, username string, organizationID uuid.UUID) (*entity.Employee, error) {
	employee, err := s.employeeRepo.GetByUsernameAndOrganizationID(ctx, username, organizationID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrEmployeeUnauthorized
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return employee, nil
}

// Create
func (s *tenderV1) Create(ctx context.Context, username string, tender entity.Tender) (*entity.Tender, error) {
	err := tender.Validate()
	if err != nil {
		return nil, NewTypedError("tender data is invalid", ErrorTypeInvalid, err)
	}

	employee, err := s.getEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	tender.Status = entity.TenderCreated
	tender.Version = 1
	tender.CreatorID = employee.ID

	createdTender, err := s.tenderRepo.Create(ctx, tender)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return createdTender, nil
}

func (s *tenderV1) GetStatus(ctx context.Context, username string, tenderID uuid.UUID) (*entity.TenderStatus, error) {
	tender, err := s.getByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	if tender.Status != entity.TenderPublished {
		_, err = s.getEmployee(ctx, username, tender.OrganizationID)
		if err != nil {
			return nil, err
		}
	}

	return &tender.Status, nil
}

func (s *tenderV1) getLimit(limit int) (int, error) {
	if limit < 0 || limit > TenderLimitMax {
		return 0, ErrTenderLimit
	}

	if limit == 0 {
		return TenderLimitDefault, nil
	}

	return limit, nil
}

// GetByServiceType
func (s *tenderV1) GetByServiceType(ctx context.Context, serviceType *entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error) {
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	if serviceType != nil {
		err := serviceType.Validate()
		if err != nil {
			return nil, NewTypedError("tender type is invalid", ErrorTypeInvalid, err)
		}
	}

	tenders, err := s.tenderRepo.GetByServiceType(ctx, serviceType, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tenders, nil
}

// GetByCreatorUsername
func (s *tenderV1) GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error) {
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	employee, err := s.getUser(ctx, username)
	if err != nil {
		return nil, err
	}

	tenders, err := s.tenderRepo.GetByCreatorID(ctx, employee.ID, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tenders, nil
}

// Update
func (s *tenderV1) Update(ctx context.Context, username string, tenderID uuid.UUID, data entity.TenderData) (*entity.Tender, error) {
	if err := data.Validate(); err != nil {
		return nil, NewTypedError("tender data is invalid", ErrorTypeInvalid, err)
	}

	tender, err := s.getByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	_, err = s.getEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	tender, err = s.tenderRepo.Update(ctx, tender.ID, data)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}

// UpdateStatus
func (s *tenderV1) UpdateStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error) {
	if err := status.Validate(); err != nil {
		return nil, NewTypedError("tender status is invalid", ErrorTypeInvalid, err)
	}

	tender, err := s.getByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	_, err = s.getEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	tender, err = s.tenderRepo.UpdateStatus(ctx, tender.ID, status)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}

// Rollback
func (s *tenderV1) Rollback(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error) {
	if version < 1 {
		return nil, ErrTenderVersion
	}

	tender, err := s.getByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	_, err = s.getEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	tender, err = s.tenderRepo.Rollback(ctx, tenderID, version)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}
