package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type BidData struct {
	Name        *string
	Description *string
}

type Bid interface {
	Create(ctx context.Context, bid entity.Bid) (*entity.Bid, error)
	GetByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	GetByTenderID(ctx context.Context, tenderID uuid.UUID, limit int, offset int) ([]entity.Bid, error)
	Update(ctx context.Context, bidID uuid.UUID, data BidData) (*entity.Bid, error)
	UpdateStatus(ctx context.Context, bidID uuid.UUID, status entity.BidStatus) (*entity.Bid, error)
	Rollback(ctx context.Context, bidID uuid.UUID, version int) (*entity.Bid, error)
}
