package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/service"
	"github.com/google/uuid"
)

type TenderReq struct {
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Status          entity.TenderStatus      `json:"status"`
	ServiceType     entity.TenderServiceType `json:"serviceType"`
	OrganizationID  uuid.UUID                `json:"organizationId"`
	CreatorUsername string                   `json:"creatorUsername"`
}

func (r TenderReq) ToTender() entity.Tender {
	return entity.Tender{
		Name:           r.Name,
		Description:    r.Description,
		Status:         r.Status,
		ServiceType:    r.ServiceType,
		OrganizationID: r.OrganizationID,
	}
}

type TenderResp struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Status      entity.TenderStatus      `json:"status"`
	ServiceType entity.TenderServiceType `json:"serviceType"`
	Version     int                      `json:"version"`
	CreatedAt   time.Time                `json:"createdAt"`
}

func (r *TenderResp) FromTender(tender *entity.Tender) {
	r.ID = tender.ID
	r.Name = tender.Name
	r.Description = tender.Description
	r.Status = tender.Status
	r.ServiceType = tender.ServiceType
	r.Version = tender.Version
	r.CreatedAt = tender.CreatedAt
}

type TendersResp []TenderResp

func (r *TendersResp) FromTenders(tenders []entity.Tender) {
	*r = make([]TenderResp, len(tenders))
	for i, tender := range tenders {
		(*r)[i].FromTender(&tender)
	}
}

// TenderGetByServiceType
// GET /tenders.
type TenderGetByServiceType struct {
	Service service.Tender
}

func (h TenderGetByServiceType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query.
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	serviceTypes := make([]entity.TenderServiceType, len(query["service_type"]))
	for i, serviceType := range query["service_type"] {
		serviceTypes[i] = entity.TenderServiceType(serviceType)
	}

	// Execute service method.
	tenders, err := h.Service.GetByServiceType(r.Context(), serviceTypes, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TendersResp
	resp.FromTenders(tenders)
	WriteValue(w, http.StatusOK, resp)
}

// TenderCreate
// POST /tenders/new.
type TenderCreate struct {
	Service service.Tender
}

func (h TenderCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request body.
	var req TenderReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		WriteReason(w, http.StatusBadRequest, err.Error())
		return
	}

	// Execute service method.
	tender, err := h.Service.Create(r.Context(), req.CreatorUsername, req.ToTender())
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TenderResp
	resp.FromTender(tender)
	WriteValue(w, http.StatusOK, resp)
}

// TenderGetByCreator
// GET /tenders/my.
type TenderGetByCreator struct {
	Service service.Tender
}

func (h TenderGetByCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query.
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	username := query.Get("username")

	// Execute service method.
	tenders, err := h.Service.GetByCreatorUsername(r.Context(), username, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TendersResp
	resp.FromTenders(tenders)
	WriteValue(w, http.StatusOK, resp)
}

// TenderGetStatus
// GET /tenders/{tenderId}/status.
type TenderGetStatus struct {
	Service service.Tender
}

func (h TenderGetStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Execute service method.
	status, err := h.Service.GetStatus(r.Context(), username, tenderID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	WriteValue(w, http.StatusOK, status)
}

// TenderUpdateStatus
// PUT /tenders/{tenderId}/status.
type TenderUpdateStatus struct {
	Service service.Tender
}

func (h TenderUpdateStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	username := query.Get("username")
	status := entity.TenderStatus(query.Get("status"))
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Execute service method.
	tender, err := h.Service.UpdateStatus(r.Context(), username, tenderID, status)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TenderResp
	resp.FromTender(tender)
	WriteValue(w, http.StatusOK, resp)
}

// TenderUpdate
// PATCH /tenders/{tenderId}/edit.
type TenderUpdate struct {
	Service service.Tender
}

func (h TenderUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Parse request body.
	var data entity.TenderData
	d := json.NewDecoder(r.Body)
	if err = d.Decode(&data); err != nil {
		WriteReason(w, http.StatusBadRequest, err.Error())
		return
	}

	// Execute service method.
	tender, err := h.Service.Update(r.Context(), username, tenderID, data)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TenderResp
	resp.FromTender(tender)
	WriteValue(w, http.StatusOK, resp)
}

// TenderRollback
// PUT /tenders/{tenderId}/rollback/{version}.
type TenderRollback struct {
	Service service.Tender
}

func (h TenderRollback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	version, _ := strconv.Atoi(r.PathValue("version"))
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Execute service method.
	tender, err := h.Service.Rollback(r.Context(), username, tenderID, version)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp TenderResp
	resp.FromTender(tender)
	WriteValue(w, http.StatusOK, resp)
}
