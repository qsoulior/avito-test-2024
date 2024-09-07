package entity

import "time"

type TenderStatus string

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

type TenderServiceType string

const (
	TenderConstruction TenderServiceType = "Construction"
	TenderDelivery     TenderServiceType = "Delivery"
	TenderManufacture  TenderServiceType = "Manufacture"
)

type Tender struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	ServiceType     string       `json:"serviceType"`
	Status          TenderStatus `json:"status"`
	OrganizationID  string       `json:"organizationId"`
	CreatorUsername string       `json:"creatorUsername"`
	Version         int          `json:"version"`
	CreatedAt       time.Time    `json:"createdAt"`
}
