package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type Bid interface {
	HasByCreatorID(ctx context.Context, creatorID uuid.UUID, tenderID uuid.UUID) error

	Create(ctx context.Context, bid entity.Bid) (*entity.Bid, error)
	GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	GetByTenderID(ctx context.Context, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	Update(ctx context.Context, bidID uuid.UUID, data entity.BidData) (*entity.Bid, error)
	UpdateStatus(ctx context.Context, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error)
	Rollback(ctx context.Context, bidID uuid.UUID, version int) (*entity.Bid, error)
}

type BidReview interface {
	Create(ctx context.Context, review entity.BidReview) (*entity.BidReview, error)
	GetByBidCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.BidReview, error)
}

type BidDecision interface {
	Create(ctx context.Context, decision entity.BidDecision) (*entity.BidDecision, error)
	GetByBidID(ctx context.Context, bidID uuid.UUID, organizationID uuid.UUID, decisionType *entity.BidStatus) ([]entity.BidDecision, error)
}
