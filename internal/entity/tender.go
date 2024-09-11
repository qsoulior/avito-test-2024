package entity

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

// TenderServiceType
type TenderServiceType string

func (t TenderServiceType) Validate() error {
	if !slices.Contains(TenderServiceTypes, t) {
		return fmt.Errorf("tender service type must be one of: %v", TenderServiceTypes)
	}
	return nil
}

const (
	TenderConstruction TenderServiceType = "Construction"
	TenderDelivery     TenderServiceType = "Delivery"
	TenderManufacture  TenderServiceType = "Manufacture"
)

var TenderServiceTypes = []TenderServiceType{TenderConstruction, TenderDelivery, TenderManufacture}

// TenderStatus
type TenderStatus string

func (s TenderStatus) Validate() error {
	if !slices.Contains(TenderStatuses, s) {
		return fmt.Errorf("tender status must be one of: %v", TenderStatuses)
	}
	return nil
}

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

var TenderStatuses = []TenderStatus{TenderCreated, TenderPublished, TenderClosed}

// Tender
type Tender struct {
	ID             uuid.UUID
	Name           string
	Description    string
	ServiceType    TenderServiceType
	Status         TenderStatus
	OrganizationID uuid.UUID
	CreatorID      uuid.UUID
	Version        int
	CreatedAt      time.Time
}

func (t Tender) Validate() error {
	if len(t.Name) > TenderNameLength {
		return ErrTenderName
	}

	if len(t.Description) > TenderDescriptionLength {
		return ErrTenderDescription
	}

	if err := t.ServiceType.Validate(); err != nil {
		return err
	}

	return t.Status.Validate()
}

const (
	TenderNameLength        = 100
	TenderDescriptionLength = 500
)

var (
	ErrTenderName        = fmt.Errorf("tender name is too long (max %d)", TenderNameLength)
	ErrTenderDescription = fmt.Errorf("tender description is too long (max %d)", TenderDescriptionLength)
)

// TenderData
type TenderData struct {
	Name        *string
	Description *string
	ServiceType *TenderServiceType
}

func (d TenderData) Validate() error {
	if d.Name != nil && len(*d.Name) > TenderNameLength {
		return ErrTenderName
	}

	if d.Description != nil && len(*d.Description) > TenderDescriptionLength {
		return ErrTenderDescription
	}

	if d.ServiceType != nil {
		return d.ServiceType.Validate()
	}

	return nil
}
