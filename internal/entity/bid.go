package entity

import (
	"time"

	"github.com/google/uuid"
)

type BidStatus string

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Canceled"
	BidApproved  BidStatus = "Approved"
	BidRejected  BidStatus = "Rejected"
)

type BidAuthorType string

const (
	BidOrganization BidAuthorType = "Organization"
	BidUser         BidAuthorType = "User"
)

type Bid struct {
	ID              uuid.UUID     `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Status          string        `json:"status"`
	TenderID        uuid.UUID     `json:"tenderId"`
	AuthorType      BidAuthorType `json:"authorType"`
	AuthorID        uuid.UUID     `json:"authorId"`
	CreatorUsername string        `json:"creatorUsername"`
	Version         int           `json:"version"`
	CreatedAt       time.Time     `json:"createdAt"`
}
