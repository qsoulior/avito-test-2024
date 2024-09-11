package http

import (
	"log/slog"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/service"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/transport/http/handler"
)

func NewMux(tenderService service.Tender, bidService service.Bid, logger *slog.Logger) http.Handler {
	if logger == nil {
		return nil
	}

	router := http.NewServeMux()
	router.Handle("GET /api/ping", handler.Ping{})

	router.Handle("GET /api/tenders", handler.TenderGetByServiceType{Service: tenderService})
	router.Handle("POST /api/tenders/new", handler.TenderCreate{Service: tenderService})
	router.Handle("GET /api/tenders/my", handler.TenderGetByCreator{Service: tenderService})
	router.Handle("GET /api/tenders/{tenderId}/status", handler.TenderGetStatus{Service: tenderService})
	router.Handle("PUT /api/tenders/{tenderId}/status", handler.TenderUpdateStatus{Service: tenderService})
	router.Handle("PATCH /api/tenders/{tenderId}/edit", handler.TenderUpdate{Service: tenderService})
	router.Handle("PUT /api/tenders/{tenderId}/rollback/{version}", handler.TenderRollback{Service: tenderService})

	router.Handle("POST /api/bids/new", handler.BidCreate{Service: bidService})
	router.Handle("GET /api/bids/my", handler.BidGetByCreator{Service: bidService})
	router.Handle("GET /api/bids/{tenderId}/list", handler.BidGetByTender{Service: bidService})
	router.Handle("GET /api/bids/{bidId}/status", handler.BidGetStatus{Service: bidService})
	router.Handle("PUT /api/bids/{bidId}/status", handler.BidUpdateStatus{Service: bidService})
	router.Handle("PATCH /api/bids/{bidId}/edit", handler.BidUpdate{Service: bidService})
	router.Handle("PUT /api/bids/{bidId}/submit_decision", handler.BidSubmitDecision{Service: bidService})
	router.Handle("PUT /api/bids/{bidId}/rollback/{version}", handler.BidRollback{Service: bidService})

	var mux http.Handler = router
	middlewares := []Middleware{RecovererMiddleware(logger), LoggerMiddleware(logger)}
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}

	return mux
}
