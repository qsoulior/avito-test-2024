package service

import (
	"context"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

const (
	BidLimitMax     = 100
	BidLimitDefault = 5
)

var (
	ErrBidNotExist        = NewTypedError("bid does not exist", ErrorTypeNotExist, nil)
	ErrBidVersionNotExist = NewTypedError("bid version does not exist", ErrorTypeNotExist, nil)
	ErrBidLimit           = NewTypedError(
		fmt.Sprintf("bid limit must be > 0 and <= %d", BidLimitMax), ErrorTypeInvalid, nil,
	)
	ErrBidOffset       = NewTypedError("bid offset must be >= 0", ErrorTypeInvalid, nil)
	ErrBidVersion      = NewTypedError("bid version must be greater than 0", ErrorTypeInvalid, nil)
	ErrBidCreator      = NewTypedError("user is not a creator", ErrorTypeForbidden, nil)
	ErrBidNotPublished = NewTypedError("bid is not published", ErrorTypeInvalid, nil)
	ErrBidCannotUpdate = NewTypedError("cannot update approved or rejected bid", ErrorTypeInvalid, nil)
)

type Bid interface {
	GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	HasByCreatorAndTender(ctx context.Context, creatorID uuid.UUID, tenderID uuid.UUID) error

	Create(ctx context.Context, username string, bid entity.Bid) (*entity.Bid, error)
	GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error)
	GetByTenderID(ctx context.Context, username string, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	GetStatus(ctx context.Context, username string, bidID uuid.UUID) (*entity.BidStatus, error)
	UpdateStatus(ctx context.Context, username string, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error)
	Update(ctx context.Context, username string, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error)
	SubmitDecision(ctx context.Context, username string, bidID uuid.UUID, decision entity.BidStatus) (*entity.Bid, error)
	Rollback(ctx context.Context, username string, bidID uuid.UUID, version int) (*entity.Bid, error)
}

const (
	BidReviewLimitMax     = 100
	BidReviewLimitDefault = 5
)

var (
	ErrBidReviewLimit = NewTypedError(
		fmt.Sprintf("bid review limit must be > 0 and <= %d", BidReviewLimitMax), ErrorTypeInvalid, nil,
	)
	ErrBidReviewOffset    = NewTypedError("bid review offset must be >= 0", ErrorTypeInvalid, nil)
	ErrBidCreatorNotExist = NewTypedError("bid creator does not exist", ErrorTypeNotExist, nil)
)

type BidReview interface {
	Create(ctx context.Context, username string, bidID uuid.UUID, description string) (*entity.BidReview, error)
	GetByBidCreator(ctx context.Context,
		requesterUsername string, creatorUsername string, tenderID uuid.UUID,
		limit int, offset int) ([]entity.BidReview, error)
}
