package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationType string

const (
	OrganizationIE  OrganizationType = "IE"
	OrganizationLLC OrganizationType = "LLC"
	OrganizationJSC OrganizationType = "JSC"
)

var OrganizationTypes = []OrganizationType{OrganizationIE, OrganizationLLC, OrganizationJSC}

type Organization struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Type        *OrganizationType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
