package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"github.com/google/uuid"
)

type tenderV1 struct {
	tenderRepo      repo.Tender
	employeeService Employee
}

func NewTenderV1(tenderRepo repo.Tender, employeeService Employee) Tender {
	if tenderRepo == nil || employeeService == nil {
		return nil
	}
	return &tenderV1{tenderRepo, employeeService}
}

// GetByID
func (s *tenderV1) GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	tender, err := s.tenderRepo.GetByID(ctx, tenderID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrTenderNotExist
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return tender, nil
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
	// Validate limit.
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	// Validate tender service type.
	if serviceType != nil {
		err := serviceType.Validate()
		if err != nil {
			return nil, NewTypedError("tender type is invalid", ErrorTypeInvalid, err)
		}
	}

	// Get published tenders by service type.
	tenders, err := s.tenderRepo.GetByServiceType(ctx, serviceType, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tenders, nil
}

// Create
func (s *tenderV1) Create(ctx context.Context, username string, tender entity.Tender) (*entity.Tender, error) {
	// Validate tender data.
	err := tender.Validate()
	if err != nil {
		return nil, NewTypedError("tender data is invalid", ErrorTypeInvalid, err)
	}

	// Get employee associated with organization.
	employee, err := s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Set tender initial values.
	tender.Status = entity.TenderCreated
	tender.Version = 1
	tender.CreatorID = employee.ID

	// Create tender.
	createdTender, err := s.tenderRepo.Create(ctx, tender)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return createdTender, nil
}

// GetByCreatorUsername
func (s *tenderV1) GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error) {
	// Validate limit.
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	// Get creator not associated with organization.
	creator, err := s.employeeService.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get tenders by creator id.
	tenders, err := s.tenderRepo.GetByCreatorID(ctx, creator.ID, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tenders, nil
}

// GetStatus
func (s *tenderV1) GetStatus(ctx context.Context, username string, tenderID uuid.UUID) (*entity.TenderStatus, error) {
	// Verify user not associated with organization.
	_, err := s.employeeService.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get tender by id.
	tender, err := s.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	return &tender.Status, nil
}

// UpdateStatus
func (s *tenderV1) UpdateStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error) {
	// Validate tender status.
	if err := status.Validate(); err != nil {
		return nil, NewTypedError("tender status is invalid", ErrorTypeInvalid, err)
	}

	// Get tender by id.
	tender, err := s.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Update tender status.
	tender, err = s.tenderRepo.UpdateStatus(ctx, tender.ID, status)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}

// Update
func (s *tenderV1) Update(ctx context.Context, username string, tenderID uuid.UUID, data entity.TenderData) (*entity.Tender, error) {
	// Validate tender data.
	if err := data.Validate(); err != nil {
		return nil, NewTypedError("tender data is invalid", ErrorTypeInvalid, err)
	}

	// Get tender by id.
	tender, err := s.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Update tender data.
	tender, err = s.tenderRepo.Update(ctx, tender.ID, data)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}

// Rollback
func (s *tenderV1) Rollback(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error) {
	// Validate tender version.
	if version < 1 {
		return nil, ErrTenderVersion
	}

	// Get tender by id.
	tender, err := s.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Rollback tender by id and version.
	tender, err = s.tenderRepo.Rollback(ctx, tender.ID, version)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return tender, nil
}
