package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/repo"
	"github.com/google/uuid"
)

// bidV1.
type bidV1 struct {
	bidRepo         repo.Bid
	decisionRepo    repo.BidDecision
	tenderService   Tender
	employeeService Employee
}

func NewBidV1(bidRepo repo.Bid, decisionRepo repo.BidDecision, tenderService Tender, employeeService Employee) Bid {
	if bidRepo == nil || tenderService == nil || employeeService == nil {
		return nil
	}
	return &bidV1{bidRepo, decisionRepo, tenderService, employeeService}
}

// GetByID.
func (s *bidV1) GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {
	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return nil, ErrBidNotExist
		}
		return nil, NewTypedError("bidRepo.GetByID", ErrorTypeInternal, err)
	}
	return bid, nil
}

// HasByCreatorAndTender.
func (s *bidV1) HasByCreatorAndTender(ctx context.Context, creatorID uuid.UUID, tenderID uuid.UUID) error {
	err := s.bidRepo.HasByCreatorID(ctx, creatorID, tenderID)
	if err != nil {
		if errors.Is(err, repo.ErrNoRows) {
			return NewTypedError("creator does not have bid for tender", ErrorTypeNotExist, nil)
		}
		return NewTypedError("bidRepo.HasByCreatorID", ErrorTypeInternal, err)
	}

	return nil
}

// Create.
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
		return nil, ErrTenderNotPublished
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

	// Set bid version.
	bid.Version = 1

	// Create bid.
	createdBid, err := s.bidRepo.Create(ctx, bid)
	if err != nil {
		return nil, NewTypedError("bidRepo.Create", ErrorTypeInternal, err)
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

// GetByCreatorUsername.
func (s *bidV1) GetByCreatorUsername(ctx context.Context,
	username string, limit int, offset int) ([]entity.Bid, error) {
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
		return nil, NewTypedError("bidRepo.GetByCreatorID", ErrorTypeInternal, err)
	}

	return bids, nil
}

// GetByTenderID.
func (s *bidV1) GetByTenderID(ctx context.Context,
	username string, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error) {
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
		return nil, NewTypedError("bidRepo.GetByTenderID", ErrorTypeInternal, err)
	}

	return bids, nil
}

// GetStatus.
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

// UpdateStatus.
func (s *bidV1) UpdateStatus(ctx context.Context,
	username string, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error) {
	// Validate bid status.
	if err := status.Validate(); err != nil {
		return nil, NewTypedError("bid status is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify bid status.
	if bid.Status == entity.BidApproved || bid.Status == entity.BidRejected {
		return nil, ErrBidCannotUpdate
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
		return nil, NewTypedError("bidRepo.UpdateStatus", ErrorTypeInternal, err)
	}

	return bid, nil
}

// Update.
func (s *bidV1) Update(ctx context.Context,
	username string, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error) {
	// Validate bid data.
	if err := data.Validate(); err != nil {
		return nil, NewTypedError("bid data is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify bid status.
	if bid.Status == entity.BidApproved || bid.Status == entity.BidRejected {
		return nil, ErrBidCannotUpdate
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
		return nil, NewTypedError("bidRepo.Update", ErrorTypeInternal, err)
	}

	return bid, nil
}

// SubmitDecision.
func (s *bidV1) SubmitDecision(ctx context.Context,
	username string, bidID uuid.UUID, decisionType entity.BidStatus) (*entity.Bid, error) {
	// Validate bid decision type.
	if err := decisionType.ValidateDesicion(); err != nil {
		return nil, NewTypedError("bid decision is invalid", ErrorTypeInvalid, err)
	}

	// Get bid by id.
	bid, err := s.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Verify bid status.
	if bid.Status != entity.BidPublished {
		return nil, ErrBidNotPublished
	}

	// Get tender by id.
	tender, err := s.tenderService.GetByID(ctx, bid.TenderID)
	if err != nil {
		return nil, err
	}

	// Verify tender status.
	if tender.Status != entity.TenderPublished {
		return nil, ErrTenderNotPublished
	}

	// Verify employee associated with organization.
	employee, err := s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	if decisionType != entity.BidRejected {
		// Get employees by organization.
		employees, err := s.employeeService.GetByOrganization(ctx, tender.OrganizationID)
		if err != nil {
			return nil, err
		}

		// Create bid decision.
		_, err = s.decisionRepo.Create(ctx, entity.BidDecision{
			BidID:          bid.ID,
			Type:           decisionType,
			OrganizationID: tender.OrganizationID,
			CreatorID:      employee.ID,
		})
		if err != nil {
			return nil, NewTypedError("decisionRepo.Create", ErrorTypeInternal, err)
		}

		// Get bid decisions.
		decisions, err := s.decisionRepo.GetByBidID(ctx, bid.ID, tender.OrganizationID, &decisionType)
		if err != nil {
			return nil, NewTypedError("decisionRepo.GetByBidID", ErrorTypeInternal, err)
		}

		// If number of decisions is less than quorum,
		// then do not update status.
		if len(decisions) < min(3, len(employees)) {
			return bid, nil
		}
	}

	// Update bid status.
	bid, err = s.bidRepo.UpdateStatus(ctx, bid.ID, decisionType)
	if err != nil {
		return nil, NewTypedError("bidRepo.UpdateStatus", ErrorTypeInternal, err)
	}

	// Update tender status.
	if decisionType == entity.BidApproved {
		_, err = s.tenderService.UpdateStatus(ctx, username, tender.ID, entity.TenderClosed)
		if err != nil {
			return nil, err
		}
	}

	return bid, nil
}

// Rollback.
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
		return nil, NewTypedError("bidRepo.Rollback", ErrorTypeInternal, err)
	}

	return bid, nil
}

// bidReviewV1.
type bidReviewV1 struct {
	reviewRepo      repo.BidReview
	bidService      Bid
	tenderService   Tender
	employeeService Employee
}

func NewBidReviewV1(
	reviewRepo repo.BidReview, bidService Bid, tenderService Tender, employeeService Employee) BidReview {
	if reviewRepo == nil || bidService == nil || tenderService == nil || employeeService == nil {
		return nil
	}
	return &bidReviewV1{reviewRepo, bidService, tenderService, employeeService}
}

// Create.
func (s *bidReviewV1) Create(ctx context.Context,
	username string, bidID uuid.UUID, description string) (*entity.BidReview, error) {
	// Get bid by id.
	bid, err := s.bidService.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Get tender by id.
	tender, err := s.tenderService.GetByID(ctx, bid.TenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	employee, err := s.employeeService.GetEmployee(ctx, username, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Set review.
	review := entity.BidReview{
		Description:    description,
		BidID:          bidID,
		OrganizationID: tender.OrganizationID,
		CreatorID:      employee.ID,
	}

	if err = review.Validate(); err != nil {
		return nil, NewTypedError("bid review data is invalid", ErrorTypeInvalid, err)
	}

	// Create review.
	createdReview, err := s.reviewRepo.Create(ctx, review)
	if err != nil {
		return nil, NewTypedError("reviewRepo.Create", ErrorTypeInternal, err)
	}

	return createdReview, nil
}

func (s *bidReviewV1) getLimit(limit int) (int, error) {
	if limit < 0 || limit > BidLimitMax {
		return 0, ErrBidReviewLimit
	}

	if limit == 0 {
		return BidLimitDefault, nil
	}

	return limit, nil
}

// GetByBidCreator.
func (s *bidReviewV1) GetByBidCreator(ctx context.Context,
	requesterUsername string, creatorUsername string, tenderID uuid.UUID,
	limit int, offset int) ([]entity.BidReview, error) {
	// Validate limit.
	limit, err := s.getLimit(limit)
	if err != nil {
		return nil, err
	}

	// Validate offset.
	if offset < 0 {
		return nil, ErrBidReviewOffset
	}

	// Get tender by id.
	tender, err := s.tenderService.GetByID(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// Verify employee associated with organization.
	_, err = s.employeeService.GetEmployee(ctx, requesterUsername, tender.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Check if creator exists.
	creator, err := s.employeeService.GetUser(ctx, creatorUsername)
	if err != nil {
		if errors.Is(err, ErrEmployeeUnauthorized) {
			return nil, NewTypedError("creator not found", ErrorTypeNotExist, nil)
		}
		return nil, err
	}

	// Check if creator have bid for tender.
	err = s.bidService.HasByCreatorAndTender(ctx, creator.ID, tender.ID)
	if err != nil {
		return nil, err
	}

	// Get bid reviews by creator.
	reviews, err := s.reviewRepo.GetByBidCreatorID(ctx, creator.ID, limit, offset)
	if err != nil {
		return nil, NewTypedError("reviewRepo.GetByBidCreatorID", ErrorTypeInternal, err)
	}

	return reviews, nil
}
