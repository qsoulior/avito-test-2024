package entity

type TenderStatus string

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

type Tender struct {
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	ServiceType     string       `json:"serviceType"`
	Status          TenderStatus `json:"status"`
	OrganizationID  string       `json:"organizationId"`
	CreatorUsername string       `json:"creatorUsername"`
}
