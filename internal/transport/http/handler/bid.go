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

type BidReq struct {
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Status          entity.BidStatus `json:"status"`
	TenderID        uuid.UUID        `json:"tenderId"`
	OrganizationID  *uuid.UUID       `json:"organizationId"`
	CreatorUsername string           `json:"creatorUsername"`
}

func (r BidReq) ToBid() entity.Bid {
	return entity.Bid{
		Name:           r.Name,
		Description:    r.Description,
		Status:         r.Status,
		TenderID:       r.TenderID,
		OrganizationID: r.OrganizationID,
	}
}

type BidResp struct {
	ID         uuid.UUID            `json:"id"`
	Name       string               `json:"name"`
	Status     entity.BidStatus     `json:"status"`
	AuthorType entity.BidAuthorType `json:"authorType"`
	AuthorID   uuid.UUID            `json:"authorId"`
	Version    int                  `json:"version"`
	CreatedAt  time.Time            `json:"createdAt"`
}

func (r *BidResp) FromBid(bid *entity.Bid) {
	r.ID = bid.ID
	r.Name = bid.Name
	r.Status = bid.Status
	if bid.OrganizationID != nil {
		r.AuthorType = entity.BidOrganization
		r.AuthorID = *bid.OrganizationID
	} else {
		r.AuthorType = entity.BidUser
		r.AuthorID = bid.CreatorID
	}
	r.Version = bid.Version
	r.CreatedAt = bid.CreatedAt
}

type BidsResp []BidResp

func (r *BidsResp) FromBids(bids []entity.Bid) {
	*r = make([]BidResp, len(bids))
	for i, bid := range bids {
		(*r)[i].FromBid(&bid)
	}
}

// BidCreate
// POST /bids/new.
type BidCreate struct {
	Service service.Bid
}

func (h BidCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request body.
	var req BidReq
	d := json.NewDecoder(r.Body)
	err := d.Decode(&req)
	if err != nil {
		WriteReason(w, http.StatusBadRequest, err.Error())
		return
	}

	// Execute service method.
	bid, err := h.Service.Create(r.Context(), req.CreatorUsername, req.ToBid())
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

// BidGetByCreator
// GET /bids/my.
type BidGetByCreator struct {
	Service service.Bid
}

func (h BidGetByCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query.
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	username := query.Get("username")

	// Execute service method.
	bids, err := h.Service.GetByCreatorUsername(r.Context(), username, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidsResp
	resp.FromBids(bids)
	WriteValue(w, http.StatusOK, resp)
}

// BidGetByTender
// GET /bids/{tenderId}/list.
type BidGetByTender struct {
	Service service.Bid
}

func (h BidGetByTender) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	username := query.Get("username")
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Execute service method.
	bids, err := h.Service.GetByTenderID(r.Context(), username, tenderID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidsResp
	resp.FromBids(bids)
	WriteValue(w, http.StatusOK, resp)
}

// BidGetStatus
// GET /bids/{bidId}/status.
type BidGetStatus struct {
	Service service.Bid
}

func (h BidGetStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Execute service method.
	status, err := h.Service.GetStatus(r.Context(), username, bidID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	WriteValue(w, http.StatusOK, status)
}

// BidUpdateStatus
// PUT /bids/{bidId}/status.
type BidUpdateStatus struct {
	Service service.Bid
}

func (h BidUpdateStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	username := query.Get("username")
	status := entity.BidStatus(query.Get("status"))
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Execute service method.
	bid, err := h.Service.UpdateStatus(r.Context(), username, bidID, status)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

// BidUpdate
// PATCH /bids/{bidId}/edit.
type BidUpdate struct {
	Service service.Bid
}

func (h BidUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Parse request body.
	var data entity.BidData
	d := json.NewDecoder(r.Body)
	if err = d.Decode(&data); err != nil {
		WriteReason(w, http.StatusBadRequest, err.Error())
		return
	}

	// Execute service method.
	bid, err := h.Service.Update(r.Context(), username, bidID, data)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

// BidSubmitDecision
// PUT /bids/{bidId}/submit_decision.
type BidSubmitDecision struct {
	Service service.Bid
}

func (h BidSubmitDecision) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	username := query.Get("username")
	decision := entity.BidStatus(query.Get("decision"))
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Execute service method.
	bid, err := h.Service.SubmitDecision(r.Context(), username, bidID, decision)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

// BidRollback
// PUT /bids/{bidId}/rollback/{version}.
type BidRollback struct {
	Service service.Bid
}

func (h BidRollback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	username := r.URL.Query().Get("username")
	version, _ := strconv.Atoi(r.PathValue("version"))
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Execute service method.
	bid, err := h.Service.Rollback(r.Context(), username, bidID, version)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

type BidReviewResp struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (r *BidReviewResp) FromBidReview(review *entity.BidReview) {
	r.ID = review.ID
	r.Description = review.Description
	r.CreatedAt = review.CreatedAt
}

type BidReviewsResp []BidReviewResp

func (r *BidReviewsResp) FromBidReviews(reviews []entity.BidReview) {
	*r = make([]BidReviewResp, len(reviews))
	for i, bid := range reviews {
		(*r)[i].FromBidReview(&bid)
	}
}

// BidReviewCreate
// PUT /bids/{bidId}/feedback.
type BidReviewCreate struct {
	Service service.BidReview
}

func (h BidReviewCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	description := query.Get("bidFeedback")
	username := query.Get("username")
	bidID, err := uuid.Parse(r.PathValue("bidId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("bidId: %s", err))
		return
	}

	// Execute service method.
	bid, err := h.Service.Create(r.Context(), username, bidID, description)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidResp
	resp.FromBid(bid)
	WriteValue(w, http.StatusOK, resp)
}

// BidReviewGetByBidCreator
// GET /bids/{tenderId}/reviews.
type BidReviewGetByBidCreator struct {
	Service service.BidReview
}

func (h BidReviewGetByBidCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request query and path.
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	creatorUsername := query.Get("authorUsername")
	requesterUsername := query.Get("requesterUsername")
	tenderID, err := uuid.Parse(r.PathValue("tenderId"))
	if err != nil {
		WriteReason(w, http.StatusBadRequest, fmt.Sprintf("tenderId: %s", err))
		return
	}

	// Execute service method.
	bid, err := h.Service.GetByBidCreator(r.Context(), requesterUsername, creatorUsername, tenderID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Write response.
	var resp BidReviewsResp
	resp.FromBidReviews(bid)
	WriteValue(w, http.StatusOK, resp)
}
