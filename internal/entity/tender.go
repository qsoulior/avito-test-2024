package entity

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TenderServiceType string

const (
	TenderConstruction TenderServiceType = "Construction"
	TenderDelivery     TenderServiceType = "Delivery"
	TenderManufacture  TenderServiceType = "Manufacture"
)

var TenderServiceTypes = []TenderServiceType{TenderConstruction, TenderDelivery, TenderManufacture}

type TenderStatus string

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

var TenderStatuses = []TenderStatus{TenderCreated, TenderPublished, TenderClosed}

type Tender struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ServiceType     TenderServiceType `json:"serviceType"`
	Status          TenderStatus      `json:"status"`
	OrganizationID  uuid.UUID         `json:"organizationId"`
	CreatorUsername string            `json:"creatorUsername"`
	Version         int               `json:"version"`
	CreatedAt       time.Time         `json:"createdAt"`
}

func (t Tender) Validate() error {
	if len(t.Name) > 100 {
		return errors.New("name is too long (max 100)")
	}

	if len(t.Description) > 500 {
		return errors.New("description is too long (max 500)")
	}

	if !slices.Contains(TenderServiceTypes, t.ServiceType) {
		return fmt.Errorf("service type must be one of: %v", TenderServiceTypes)
	}

	if !slices.Contains(TenderStatuses, t.Status) {
		return fmt.Errorf("status must be one of: %v", TenderStatuses)
	}

	if t.Version < 1 {
		return errors.New("version must be greater than or equal to 1")
	}

	return nil
}
