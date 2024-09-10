package service

import (
	"context"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"github.com/google/uuid"
)

const (
	TenderLimitMax     = 100
	TenderLimitDefault = 5
)

var (
	ErrTenderNotExist = NewTypedError("tender does not exist", ErrorTypeNotExist, nil)
	ErrTenderLimit    = NewTypedError(
		fmt.Sprintf("tender limit must be > 0 and <= %d", TenderLimitMax), ErrorTypeInvalid, nil,
	)
	ErrTenderVersion = NewTypedError("tender version must be greater than 0", ErrorTypeInvalid, nil)
)

type Tender interface {
	GetByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error)
	GetByServiceType(ctx context.Context, serviceTypes []entity.TenderServiceType, limit int, offset int) ([]entity.Tender, error)
	Create(ctx context.Context, username string, tender entity.Tender) (*entity.Tender, error)
	GetByCreatorUsername(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error)
	GetStatus(ctx context.Context, username string, tenderID uuid.UUID) (*entity.TenderStatus, error)
	UpdateStatus(ctx context.Context, username string, tenderID uuid.UUID, status entity.TenderStatus) (*entity.Tender, error)
	Update(ctx context.Context, username string, tenderID uuid.UUID, data entity.TenderData) (*entity.Tender, error)
	Rollback(ctx context.Context, username string, tenderID uuid.UUID, version int) (*entity.Tender, error)
}
