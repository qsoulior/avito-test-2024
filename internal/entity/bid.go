package entity

import "time"

type BidStatus string

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Canceled"
	BidApproved  BidStatus = "Approved"
	BidRejected  BidStatus = "Rejected"
)

type Bid struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	TenderID        string    `json:"tenderId"`
	OrganizationID  string    `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
	AuthorType      string    `json:"authorType"`
	AuthorID        string    `json:"authorId"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"createdAt"`
}
