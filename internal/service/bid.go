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
	ErrBidNotExist = NewTypedError("bid does not exist", ErrorTypeNotExist, nil)
	ErrBidLimit    = NewTypedError(
		fmt.Sprintf("bid limit must be > 0 and <= %d", BidLimitMax), ErrorTypeInvalid, nil,
	)
	ErrBidVersion = NewTypedError("bid version must be greater than 0", ErrorTypeInvalid, nil)
)

type Bid interface {
	GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	Create(ctx context.Context, username string, bid entity.Bid) (*entity.Bid, error)
	GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error)
	GetByTenderID(ctx context.Context, username string, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	GetStatus(ctx context.Context, username string, bidID uuid.UUID) (*entity.BidStatus, error)
	UpdateStatus(ctx context.Context, username string, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error)
	Update(ctx context.Context, username string, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error)
	SubmitDecision(ctx context.Context, username string, bidID uuid.UUID, decision entity.BidStatus) (*entity.Bid, error)
	Rollback(ctx context.Context, username string, bidID uuid.UUID, version int) (*entity.Bid, error)
}
