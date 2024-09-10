package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type Tender interface {
	Create(ctx context.Context, tender entity.Tender) (*entity.Tender, error)
	GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error)
	GetByServiceType(ctx context.Context, serviceTypes []entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int, offset int) ([]entity.Tender, error)
	Update(ctx context.Context, tenderID uuid.UUID, data entity.TenderData) (*entity.Tender, error)
	UpdateStatus(ctx context.Context, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error)
	Rollback(ctx context.Context, tenderID uuid.UUID, version int) (*entity.Tender, error)
}
