package entity

import (
	"time"

	"github.com/google/uuid"
)

type TenderServiceType string

const (
	TenderConstruction TenderServiceType = "Construction"
	TenderDelivery     TenderServiceType = "Delivery"
	TenderManufacture  TenderServiceType = "Manufacture"
)

type TenderStatus string

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

type Tender struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ServiceType     TenderServiceType `json:"serviceType"`
	Status          TenderStatus      `json:"status"`
	OrganizationID  uuid.UUID         `json:"organization_id"`
	CreatorUsername string            `json:"creator_username"`
	Version         int               `json:"version"`
	CreatedAt       time.Time         `json:"createdAt"`
}
