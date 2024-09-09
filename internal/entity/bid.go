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
	ID             uuid.UUID
	Name           string
	Description    string
	Status         BidStatus
	TenderID       uuid.UUID
	OrganizationID *uuid.UUID
	CreatorID      uuid.UUID
	Version        int
	CreatedAt      time.Time
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

	if b.Version < 1 {
		return errors.New("version must be greater than or equal to 1")
	}

	return nil
}
