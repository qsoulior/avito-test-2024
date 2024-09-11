package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"github.com/google/uuid"
)

type bidV1 struct {
	bidRepo         repo.Bid
	tenderService   Tender
	employeeService Employee
}

func NewBidV1(bidRepo repo.Bid, tenderService Tender, employeeService Employee) Bid {
	if bidRepo == nil || tenderService == nil || employeeService == nil {
		return nil
	}
	return &bidV1{bidRepo, tenderService, employeeService}
}

// GetByID
func (s *bidV1) GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {
	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrBidNotExist
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}
	return bid, nil
}

// Create
func (s *bidV1) Create(ctx context.Context, username string, bid entity.Bid) (*entity.Bid, error) {
	// Validate data to create bid.
	err := bid.Validate()
	if err != nil {
		return nil, NewTypedError("bid data is invalid", ErrorTypeInvalid, err)
	}

	// Verify tender status.
	status, err := s.tenderService.GetStatus(ctx, username, bid.TenderID)
	if err != nil {
		return nil, err
	}
	if status == nil || *status != entity.TenderPublished {
		return nil, NewTypedError("tender is not published", ErrorTypeInvalid, nil)
	}

	// Set bid employee or private user.
	if bid.OrganizationID != nil {
		employee, err := s.employeeService.GetEmployee(ctx, username, *bid.OrganizationID)
		if err != nil {
			return nil, err
		}
		bid.CreatorID = employee.ID
	} else {
		user, err := s.employeeService.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}
		bid.CreatorID = user.ID
	}

	// Set bid status and version.
	bid.Status = entity.BidCreated
	bid.Version = 1

	// Create bid.
	createdBid, err := s.bidRepo.Create(ctx, bid)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return createdBid, nil
}

func (s *bidV1) getLimit(limit int) (int, error) {
	if limit < 0 || limit > BidLimitMax {
		return 0, ErrBidLimit
	}

	if limit == 0 {
		return BidLimitDefault, nil
	}

	return limit, nil
}

// GetByCreatorUsername
func (s *bidV1) GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error) {
	// Validate limit.
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	// Validate offset.
	if offset < 0 {
		return nil, ErrBidOffset
	}

	// Get creator not associated with organization.
	employee, err := s.employeeService.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get bids by creator id.
	bids, err := s.bidRepo.GetByCreatorID(ctx, employee.ID, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bids, nil
}

// GetByTenderID
func (s *bidV1) GetByTenderID(ctx context.Context, username string, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error) {
	// Validate limit.
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	// Validate offset.
	if offset < 0 {
		return nil, ErrBidOffset
	}

	// Get tender by id.
	tender, err := s.tenderService.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Get bids by tender id.
	bids, err := s.bidRepo.GetByTenderID(ctx, tender.ID, limit, offset)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bids, nil
}

// GetStatus
func (s *bidV1) GetStatus(ctx context.Context, username string, bidID uuid.UUID) (*entity.BidStatus, error) {
	// Verify user not associated with organization.
	_, err := s.employeeService.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	return &bid.Status, nil
}

// UpdateStatus
func (s *bidV1) UpdateStatus(ctx context.Context, username string, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error) {
	// Validate bid status.
	if err := status.Validate(); err != nil {
		return nil, NewTypedError("bid status is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify employee ot private user.
	if bid.OrganizationID != nil {
		_, err = s.employeeService.GetEmployee(ctx, username, *bid.OrganizationID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := s.employeeService.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}
		if user.ID != bid.CreatorID {
			return nil, ErrBidCreator
		}
	}

	// Update bid status.
	bid, err = s.bidRepo.UpdateStatus(ctx, bid.ID, status)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bid, nil
}

// Update
func (s *bidV1) Update(ctx context.Context, username string, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error) {
	// Validate bid data.
	if err := data.Validate(); err != nil {
		return nil, NewTypedError("bid data is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify employee ot private user.
	if bid.OrganizationID != nil {
		_, err = s.employeeService.GetEmployee(ctx, username, *bid.OrganizationID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := s.employeeService.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}
		if user.ID != bid.CreatorID {
			return nil, ErrBidCreator
		}
	}

	// Update bid data.
	bid, err = s.bidRepo.Update(ctx, bid.ID, data)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bid, nil
}

// SubmitDecision
func (s *bidV1) SubmitDecision(ctx context.Context, username string, bidID uuid.UUID, decision entity.BidStatus) (*entity.Bid, error) {
	// Validate bid decision.
	if err := decision.ValidateDesicion(); err != nil {
		return nil, NewTypedError("bid decision is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Get tender by id.
	tender, err := s.tenderService.GetByID(ctx, bid.TenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Update bid status.
	bid, err = s.bidRepo.UpdateStatus(ctx, bid.ID, decision)
	if err != nil {
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bid, nil
}

// Rollback
func (s *bidV1) Rollback(ctx context.Context, username string, bidID uuid.UUID, version int) (*entity.Bid, error) {
	// Validate bid version.
	if version < 1 {
		return nil, ErrBidVersion
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify employee ot private user.
	if bid.OrganizationID != nil {
		_, err = s.employeeService.GetEmployee(ctx, username, *bid.OrganizationID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := s.employeeService.GetUser(ctx, username)
		if err != nil {
			return nil, err
		}
		if user.ID != bid.CreatorID {
			return nil, ErrBidCreator
		}
	}

	// Rollback bid by id and version.
	bid, err = s.bidRepo.Rollback(ctx, bid.ID, version)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrBidVersionNotExist
		}
		return nil, NewTypedError("", ErrorTypeInternal, err)
	}

	return bid, nil
}
