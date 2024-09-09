package repo

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type TenderData struct {
	Name        *string
	Description *string
	ServiceType *entity.TenderServiceType
}

type Tender interface {
	Create(ctx context.Context, tender entity.Tender) (*entity.Tender, error)
	GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error)
	GetByServiceType(ctx context.Context, serviceType *entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error)
	GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error)
	Update(ctx context.Context, tenderID uuid.UUID, username string, data TenderData) (*entity.Tender, error)
	UpdateStatus(ctx context.Context, tenderID uuid.UUID, username string, status entity.TenderStatus) (*entity.Tender, error)
	Rollback(ctx context.Context, tenderID uuid.UUID, username string, version int) (*entity.Tender, error)
}
