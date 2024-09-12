package entity

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

// BidStatus
type BidStatus string

func (s BidStatus) Validate() error {
	if !slices.Contains(BidStatuses, s) {
		return fmt.Errorf("bid status must be one of: %v", BidStatuses)
	}
	return nil
}

func (s BidStatus) ValidateDesicion() error {
	if !slices.Contains(BidDecisionTypes, s) {
		return fmt.Errorf("bid decision must be one of: %v", BidDecisionTypes)
	}
	return nil
}

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Canceled"
	BidApproved  BidStatus = "Approved"
	BidRejected  BidStatus = "Rejected"
)

var (
	BidStatuses      = []BidStatus{BidCreated, BidPublished, BidCanceled}
	BidDecisionTypes = []BidStatus{BidApproved, BidRejected}
)

// BidAuthorType
type BidAuthorType string

func (t BidAuthorType) Validate() error {
	if !slices.Contains(BidAuthorTypes, t) {
		return fmt.Errorf("bid author type must be one of: %v", BidAuthorTypes)
	}
	return nil
}

const (
	BidOrganization BidAuthorType = "Organization"
	BidUser         BidAuthorType = "User"
)

var BidAuthorTypes = []BidAuthorType{BidOrganization, BidUser}

// Bid
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
	if len(b.Name) > BidNameLength {
		return ErrBidName
	}

	if len(b.Description) > BidDescriptionLength {
		return ErrBidDescription
	}

	return b.Status.Validate()
}

const (
	BidNameLength              = 100
	BidDescriptionLength       = 500
	BidReviewDescriptionLength = 1000
)

var (
	ErrBidName              = fmt.Errorf("bid name is too long (max %d)", BidNameLength)
	ErrBidDescription       = fmt.Errorf("bid description is too long (max %d)", BidDescriptionLength)
	ErrBidReviewDescription = fmt.Errorf("bid review description is too long (max %d)", BidReviewDescriptionLength)
)

// BidData
type BidData struct {
	Name        *string
	Description *string
}

func (d BidData) Validate() error {
	if d.Name != nil && len(*d.Name) > BidNameLength {
		return ErrBidName
	}

	if d.Description != nil && len(*d.Description) > BidDescriptionLength {
		return ErrBidDescription
	}

	return nil
}

// BidReview
type BidReview struct {
	ID             uuid.UUID
	Description    string
	BidID          uuid.UUID
	OrganizationID uuid.UUID
	CreatorID      uuid.UUID
	CreatedAt      time.Time
}

func (r BidReview) Validate() error {
	if len(r.Description) > BidReviewDescriptionLength {
		return ErrBidReviewDescription
	}
	return nil
}

// BidDecision
type BidDecision struct {
	ID             uuid.UUID
	BidID          uuid.UUID
	Type           BidStatus
	OrganizationID uuid.UUID
	CreatorID      uuid.UUID
	CreatedAt      time.Time
}
