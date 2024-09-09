package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

type Tender interface {
	Create(ctx context.Context, username string, tender entity.Tender) (*entity.Tender, error)
	GetStatusByID(ctx context.Context, username string, tenderID uuid.UUID) (*entity.TenderStatus, error)
	GetByServiceType(ctx context.Context, serviceType *entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error)
	GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error)
	Update(ctx context.Context, username string, tenderID uuid.UUID, data entity.TenderData) (*entity.Tender, error)
	UpdateStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error)
	Rollback(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error)
}
