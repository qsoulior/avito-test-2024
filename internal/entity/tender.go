package entity

import "time"

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
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	ServiceType    string       `json:"serviceType"`
	Status         TenderStatus `json:"status"`
	OrganizationID string       `json:"organizationId"`
	Version        int          `json:"version"`
	CreatedAt      time.Time    `json:"createdAt"`
}
