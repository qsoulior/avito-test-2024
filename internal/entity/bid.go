package entity

import (
	"errors"
	"fmt"
	"slices"
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

var BidStatuses = []BidStatus{BidCreated, BidPublished, BidCanceled, BidApproved, BidRejected}

type BidAuthorType string

const (
	BidOrganization BidAuthorType = "Organization"
	BidUser         BidAuthorType = "User"
)

var BidAuthorTypes = []BidAuthorType{BidOrganization, BidUser}

type Bid struct {
	ID              uuid.UUID     `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Status          BidStatus     `json:"status"`
	TenderID        uuid.UUID     `json:"tenderId"`
	AuthorType      BidAuthorType `json:"authorType"`
	AuthorID        uuid.UUID     `json:"authorId"`
	CreatorUsername string        `json:"creatorUsername"`
	Version         int           `json:"version"`
	CreatedAt       time.Time     `json:"createdAt"`
}

func (b Bid) Validate() error {
	if len(b.Name) > 100 {
		return errors.New("name is too long (max 100)")
	}

	if len(b.Description) > 500 {
		return errors.New("description is too long (max 500)")
	}

	if !slices.Contains(BidStatuses, b.Status) {
		return fmt.Errorf("status must be one of: %v", BidStatuses)
	}

	if !slices.Contains(BidAuthorTypes, b.AuthorType) {
		return fmt.Errorf("author type must be one of: %v", BidAuthorTypes)
	}

	if b.Version < 1 {
		return errors.New("version must be greater than or equal to 1")
	}

	return nil
}
